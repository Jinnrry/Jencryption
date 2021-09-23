/**
 * 将canvas的图片数组转换成三维数组
 *
 * @param data
 * @param width
 * @param height
 * @returns {*[]}
 */
function canvasImgArg2Arg(data, width, height) {
    let ret = new Array(width)
    for (let i = 0; i < ret.length; i++) {
        ret[i] = new Array(height)
    }

    for (let i = 0; i < data.length; i += 4) {
        let x = (i / 4) % width
        let y = parseInt(parseInt(i / 4) / width)
        ret[x][y] = [data[i], data[i + 1], data[i + 2], data[i + 3],]
    }
    return ret
}

/**
 *
 * 将三维数组生成canvas使用的一维数组
 *
 * @param data
 * @param width
 * @param height
 * @returns {ImageData}
 */
function genImgObj(data, width, height) {
    let ret = new ImageData(width, height);
    for (let x = 0; x < width; x++) {
        for (let y = 0; y < height; y++) {
            ret.data[y * width * 4 + x * 4] = data[x][y][0]
            ret.data[y * width * 4 + x * 4 + 1] = data[x][y][1]
            ret.data[y * width * 4 + x * 4 + 2] = data[x][y][2]
            ret.data[y * width * 4 + x * 4 + 3] = data[x][y][3]
        }
    }
    return ret
}


function point2location(point, width, height) {
    return [point % width, parseInt(point / width)]
}

/**
 * md5后去后4位
 * @param str
 * @returns {number}
 */
function md5hash(str) {
    var hash = md5(str)
    return parseInt(hash.slice(24, 32), 16)
}

/**
 * 颜色值乱序
 * @param data
 * @param passwordHash
 * @returns {*|*[]}
 */
function shuffleColor(data, passwordHash) {
    let ret = [...data]
    for (let i = 0; i < 256; i++) {
        let tmp = ret[i]
        let idx = md5hash(passwordHash + i) % 256
        ret[i] = ret[idx]
        ret[idx] = tmp
    }
    return ret
}

/**
 * 数组乱序
 * @param data
 * @param passwordHash
 * @returns {*|*[]}
 */
function shufflePoint(data, passwordHash) {
    let ret = [...data]
    let lens = ret.length
    for (let i = 0; i < lens; i++) {
        let tmp = ret[i]
        let idx = md5hash(passwordHash + i) % lens
        ret[i] = ret[idx]
        ret[idx] = tmp
    }
    return ret
}

/**
 * 将数组的值与索引翻转
 * @param data
 * @returns {*[]}
 */
function evertArray(data) {
    var lens = data.length
    var ret = []
    for (let i = 0; i < lens; i++) {
        ret[data[i]] = i
    }
    return ret
}


function getNewPoint(pointsOffset, index, imgWidth, imgHeight) {
    let newIdx = pointsOffset[index]
    return point2location(newIdx, imgWidth, imgHeight)
}


/**
 * 解密一张图片，需要传入一张图片的三维数组
 *
 * @param img
 * @param password
 * @returns {*[]}
 * @constructor
 */
function Decrypt(img, password) {
    let imgHeight = img[0].length
    let imgWidth = img.length
    let imgSize = imgWidth * imgHeight

    let newRgba = new Array(imgWidth)
    for (let i = 0; i < newRgba.length; i++) {
        newRgba[i] = new Array(imgHeight)
    }


    // 密码hash
    let hash = md5hash(password)

    let pointOffset = [];
    for (let i = 0; i < imgSize; i++) {
        pointOffset[i] = i
    }

    let colorOffset = []
    for (let i = 0; i < 256; i++) {
        colorOffset[i] = i
    }

    let rOffset = shuffleColor(colorOffset, hash)
    let gOffset = shuffleColor(rOffset, hash)
    let bOffset = shuffleColor(gOffset, hash)

    rOffset = evertArray(rOffset)
    gOffset = evertArray(gOffset)
    bOffset = evertArray(bOffset)

    pointOffset = shufflePoint(pointOffset, hash)
    pointOffset = evertArray(pointOffset)

    for (let i = 0; i < imgSize; i++) {
        let oldP = point2location(i, imgWidth, imgHeight)
        let newP = getNewPoint(pointOffset, i, imgWidth, imgHeight)
        let c = img[oldP[0]][oldP[1]]
        let newR = c[0]
        let newG = c[1]
        let newB = c[2]
        let newA = c[3]
        // js canvas中没法获取未预乘的RGBA值，因此针对携带alpha 通道值的点不做颜色偏移
        if (newA !== 255) {
            newRgba[newP[0]][newP[1]] = [newR, newG, newB, newA]
        } else {
            newRgba[newP[0]][newP[1]] = [rOffset[newR], gOffset[newG], bOffset[newB], newA]
        }
    }

    return newRgba
}


/**
 * 加密一张图片，需要传入一个图片的三维数组
 * img[x][y][]
 *
 * @param img
 * @param password
 * @returns {*[]}
 * @constructor
 */
function Encrypt(img, password) {
    let imgHeight = img[0].length
    let imgWidth = img.length
    let imgSize = imgWidth * imgHeight

    let newRgba = new Array(imgWidth)
    for (let i = 0; i < newRgba.length; i++) {
        newRgba[i] = new Array(imgHeight)
    }

    // 密码hash
    let hash = md5hash(password)

    var pointOffset = []
    for (let i = 0; i < imgSize; i++) {
        pointOffset[i] = i
    }

    var colorOffset = []
    for (let i = 0; i < 256; i++) {
        colorOffset[i] = i
    }

    let rOffset = shuffleColor(colorOffset, hash)
    let gOffset = shuffleColor(rOffset, hash)
    let bOffset = shuffleColor(gOffset, hash)

    pointOffset = shufflePoint(pointOffset, hash)

    for (let i = 0; i < imgSize; i++) {
        let oldP = point2location(i, imgWidth, imgHeight)
        let newP = getNewPoint(pointOffset, i, imgWidth, imgHeight)
        let c = img[oldP[0]][oldP[1]]
        let newR = c[0]
        let newG = c[1]
        let newB = c[2]
        let newA = c[3]
        // js canvas中没法获取未预乘的RGBA值，因此针对携带alpha 通道值的点不做颜色偏移
        if (newA !== 255) {
            newRgba[newP[0]][newP[1]] = [newR, newG, newB, newA]
        } else {
            newRgba[newP[0]][newP[1]] = [rOffset[newR], gOffset[newG], bOffset[newB], newA]
        }

    }

    return newRgba

}

/**
 * 解密当前页面上的全部图片
 * @param password
 * @constructor
 */
function DecryptAllImage(password) {
    let canvas = document.createElement("canvas")
    let ctx = canvas.getContext('2d');
    let imgs = document.getElementsByTagName("img")
    for (let i = 0; i < imgs.length; i++) {
        let img = new Image();
        img.crossOrigin = 'anonymous';
        img.src = imgs[i].src;
        let filePath = imgs[i].src
        //获取最后一个.的位置
        let index = filePath.lastIndexOf(".");
        //获取后缀
        let ext = filePath.substr(index + 1);
        if (ext !== "png" && ext !== "jpeg" && ext !== "jpg") {
            continue
        }
        img.onload = function () {
            canvas.width = img.width
            canvas.height = img.height
            ctx.drawImage(img, 0, 0);
            let imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
            let imgArg = canvasImgArg2Arg(imageData.data, canvas.width, canvas.height)
            let eImgArg = Decrypt(imgArg, password)
            ctx.putImageData(genImgObj(eImgArg, canvas.width, canvas.height), 0, 0)
            imgs[i].src = canvas.toDataURL("image/png", 1) //获取Base64编码
            imgs[i].style.visibility = "visible"
        };
        imgs[i].style.visibility = "hidden"
    }
}

/**
 * 解密单张图片
 * demo：
 * DecryptImage(document.getElementById("img"),"123")
 *
 * @param imgElement
 * @param password
 * @constructor
 */
function DecryptImage(imgElement, password) {
    let canvas = document.createElement("canvas")
    let ctx = canvas.getContext('2d', {
        alpha: true
    });
    let img = new Image();
    img.crossOrigin = 'anonymous';
    img.src = imgElement.src;
    img.onload = function () {
        canvas.width = img.width
        canvas.height = img.height
        ctx.drawImage(img, 0, 0);
        let imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        let imgArg = canvasImgArg2Arg(imageData.data, canvas.width, canvas.height)
        let eImgArg = Decrypt(imgArg, password)
        ctx.putImageData(genImgObj(eImgArg, canvas.width, canvas.height), 0, 0)
        imgElement.src = canvas.toDataURL("image/png", 1); //获取Base64编码
    };

}


/**
 * 加密单张图片
 * demo：
 * EncryptImage(document.getElementById("img"),"123")
 *
 * @param imgElement
 * @param password
 * @constructor
 */
function EncryptImage(imgElement, password) {
    let canvas = document.createElement("canvas")
    let ctx = canvas.getContext('2d', {
        alpha: true
    });
    let img = new Image();
    img.crossOrigin = 'anonymous';
    img.src = imgElement.src;
    img.onload = function () {
        canvas.width = img.width
        canvas.height = img.height
        ctx.drawImage(img, 0, 0);
        let imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        let imgArg = canvasImgArg2Arg(imageData.data, canvas.width, canvas.height)
        let eImgArg = Encrypt(imgArg, password)
        ctx.putImageData(genImgObj(eImgArg, canvas.width, canvas.height), 0, 0)
        imgElement.src = canvas.toDataURL("image/png", 1); //获取Base64编码
    };

}