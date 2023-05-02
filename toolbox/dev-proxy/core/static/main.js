function get(url, successFuc, failFunc) {
    loading()
    // 1、创建xhr对象
    const xhr = new XMLHttpRequest();
    // 2、调用open函数   创建请求
    xhr.open("GET", url);
    // 3、调用send函数   发起请求
    xhr.send();
    // 4、监听onreadystatechange事件
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            //数据获取成功
            // console.log(xhr.responseText);
            // xhr.responseText返回的数据中有一个status和xhr.status完全不一样
            if (successFuc) {
                successFuc(xhr.responseText);
            }
        } else {
            if (failFunc) {
                failFunc(xhr.responseText);
            }
        }
        loaded()
    };
}

function postJSON(url, body) {
    //   1、创建xhr对象
    const xhr = new XMLHttpRequest();
    // 2、调用open函数 创建请求
    xhr.open("POST", url);
    // 3、设置conten-type属性  在发送之前对url所有字符编码
    //  给该请求增加额外的请求头部
    xhr.setRequestHeader(
        "Content-Type",
        "application/json"
    );
    // 4、调用send函数   发起请求
    xhr.send(body);
    // 5、监听事件
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log(xhr.responseText);
        }
    };
}


// 显示loading遮罩层
function loading() {
    let mask_bg = document.createElement('div')
    mask_bg.id = 'mask_bg'
    mask_bg.style.position = 'absolute'
    mask_bg.style.top = '0px'
    mask_bg.style.left = '0px'
    mask_bg.style.width = '100%'
    mask_bg.style.height = '100%'
    mask_bg.style.backgroundColor = '#777'
    mask_bg.style.opacity = 0.6
    mask_bg.style.zIndex = 10001
    document.body.appendChild(mask_bg)
    let mask_msg = document.createElement('div')
    mask_msg.style.position = 'absolute'
    mask_msg.style.top = '35%'
    mask_msg.style.left = '42%'
    mask_msg.style.backgroundColor = 'white'
    mask_msg.style.border = '#336699 1px solid'
    mask_msg.style.textAlign = 'center'
    mask_msg.style.fontSize = '1.1em'
    mask_msg.style.fontWeight = 'bold'
    mask_msg.style.padding = '0.5em 3em 0.5em 3em'
    mask_msg.style.zIndex = 10002
    mask_msg.innerText = 'Loading ...'
    mask_bg.appendChild(mask_msg)
}

// 关闭遮罩层
function loaded() {
    let mask_bg = document.getElementById('mask_bg')
    if (mask_bg != null) mask_bg.parentNode.removeChild(mask_bg)
}

Date.prototype.format = function (format) {
    let o = {
        "M+": this.getMonth() + 1, //month
        "d+": this.getDate(), //day
        "h+": this.getHours(), //hour
        "m+": this.getMinutes(), //minute
        "s+": this.getSeconds(), //second
        "q+": Math.floor((this.getMonth() + 3) / 3), //quarter
        "S": this.getMilliseconds() //millisecond
    }
    if (/(y+)/.test(format)) format = format.replace(RegExp.$1,
        (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (let k in o) if (new RegExp("(" + k + ")").test(format))
        format = format.replace(RegExp.$1,
            RegExp.$1.length === 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length));
    return format;
}