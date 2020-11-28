//数字键盘输入框
function NumKeyBoard(options) {
    var defaultOptions = {
        insideInput: true,  //是否显示input框
        precision: 2,       //精确度
        minVal: 0,
        maxVal: 1000
    }

    var param = $.extend(defaultOptions, options);

    this.minVal = param.minVal;
    this.maxVal = param.maxVal;

    this.precision = param.precision;
    if ($("#numKeyBoard-box").length) {
        $("#numKeyBoard-box").css("display", "block");
        return;
    }

    //创建数字键盘table
    const inputStr = '<input class="numKeyBoardInput" readonly="readonly"/>';
    let tableStr = '<table border="0" cellspacing="0" cellpadding="0" class="numKeyBoardTable">';
    tableStr += '<tr class="numTr"><td>1</td><td>2</td><td>3</td></tr>';
    tableStr += '<tr class="numTr"><td>4</td><td>5</td><td>6</td></tr>';
    tableStr += '<tr class="numTr"><td>7</td><td>8</td><td>9</td></tr>';
    tableStr += '<tr><td class="delTd">删除</td><td class="zeroTd">0</td>';
    tableStr += '<td class="dotTd">.</td></tr>';
    tableStr += '</table>';

    let nbtn = '<div class="fBottom">';
    nbtn += '<div class="saveBtn">确认</div><div class="clearBtn">清除</div><div class="cancelBtn" >取消</div>';
    nbtn += '</div>';
    const nHtml = inputStr + tableStr + nbtn;
    const nbox = "<div id='numKeyBoard-box'><div class='numPadLayerBox'><div class='numPadLayerBox-cot'>" + nHtml + "</div></div></div>";
    $('body').append(nbox);

    this.inputValue = "";
    if (param.insideInput) {
        $(".numKeyBoardInput").css("display", 'block')
    } else {
        $(".numKeyBoardInput").css("display", 'none')
    }

    this.input = $(".numKeyBoardInput");
    this.input.attr('type', 'text')
    var self = this;

    //数字按键
    $("#numKeyBoard-box .numTr td").click(function (e) {
        if (self.inputValue.length === 1 && self.inputValue === 0) return;
        if (self.inputValue) {
            if (!validPrecision(self.inputValue)) return;
        }
        var ev = e || window.event;
        var clickEl = ev.element || ev.target;
        var tdVal = $(clickEl).html();
        self.inputValue = self.inputValue + tdVal;
        self.input.val(self.inputValue)
    })

    //删除按键
    $("#numKeyBoard-box .delTd").click(function () {
        if (self.inputValue) {
            self.inputValue = self.inputValue.substr(0, self.inputValue.length - 1);
            self.input.val(self.inputValue)
        }
    })

    //小数点按键
    $("#numKeyBoard-box .dotTd").click(function () {
        if (!self.precision) return;

        if (self.inputValue) {
            if (self.inputValue.indexOf(".") !== -1) return;
            self.inputValue = self.inputValue + '.'
            self.input.val(self.inputValue)
        } else {
            self.inputValue = "0.";
            self.input.val(self.inputValue)
        }
    })

    //0按键
    $("#numKeyBoard-box .zeroTd").click(function () {
        if (self.inputValue) {
            if (!validPrecision(self.inputValue)) return;
            if (self.inputValue.length === 1 && self.inputValue === 0) return;
            self.inputValue = self.inputValue + '0';
            self.input.val(self.inputValue)
        } else {
            self.inputValue = "0";
            self.input.val(self.inputValue)
        }
    })


    //验证精度
    function validPrecision(str) {
        if (str && str.indexOf(".") !== -1) {
            let len = str.split(".")[1].length;
            if (len === self.precision) return false;
        }
        return true;
    }

    //确认
    $(".saveBtn").click(function () {
        const val = $(".numKeyBoardInput").val();

        const re = /\.$/;
        if (re.test(val)) return;

        if (self.minVal && val < self.minVal) {
            alert('输入值不可小于' + self.minVal);
            return false;
        }

        if (self.maxVal && val > self.maxVal) {
            alert('输入值不可大于' + self.maxVal);
            return false;
        }

        if (self.confirmCallback) {
            self.confirmCallback(val);
        }

        $("#numKeyBoard-box").hide();

    });


    //取消
    $(".cancelBtn").click(function () {
        $("#numKeyBoard-box").hide();
    });

    //清除
    $(".clearBtn").click(function () {
        self.setNumVal("");
    });

    //确认的callBack

    //获取input数值
    this.getNumVal = function () {
        return self.inputValue;
    };

    //设置input数值
    this.setNumVal = function (val) {
        self.input.val(val);
        self.inputValue = val;
    };

    //打开弹窗
    this.open = function (callback) {
        $("#numKeyBoard-box").show();
        self.confirmCallback = callback;
    }
}