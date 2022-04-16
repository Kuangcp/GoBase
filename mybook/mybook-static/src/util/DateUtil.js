/**
 * 仿照 moment 的接口，简单实现 formatDate
 * @param   {Number} timeStamp    待格式化的 unix 时间戳
 * @return  {String}              格式化后的日期字符串
 */

const formatter = function (timeStamp) {
    const fixLen = function (num, len) {
        let r = `${num}`;
        while (r.length < len) {
            r = `0${r}`;
        }
        return r;
    };
    const date = new Date(timeStamp);
    return {
        format(pattern) {
            const str = typeof pattern === 'string' ? pattern : 'yyyy-MM-dd';
            if (!isNaN(date.getTime())) {
                return str
                    .replace(/yyyy/i, fixLen(date.getFullYear(), 4))
                    .replace(/MM/, fixLen(date.getMonth() + 1, 2))
                    .replace(/dd/i, fixLen(date.getDate(), 2))
                    .replace(/hh/i, fixLen(date.getHours(), 2))
                    .replace(/mm/, fixLen(date.getMinutes(), 2))
                    .replace(/ss/i, fixLen(date.getSeconds(), 2));
            }
            return '';
        },
        formatDate() {
            return this.format("YYYY-MM-dd")
        }
    };
}

// TODO 自然月 前进 后退
// const fillRelativeDate = function (picker) {
// }
const day = 3600 * 1000 * 24
const passWeek = day * 7
const passMonth = day * 30
const passTriMonth = day * 90
const passHalfYear = day * 183
const passYear = day * 365
const passTriYear = day * 365 * 3
const passFiveYear = day * 365 * 5

const calDate = function (offset) {
    const end = new Date();
    const start = new Date();
    start.setTime(start.getTime() - offset);
    return [start, end]
}

const fillDate = function (picker, offset) {
    const end = new Date();
    const start = new Date();
    start.setTime(start.getTime() - offset);
    picker.$emit("pick", [start, end]);
}

const dateShortCut = [
    {
        text: "今年",
        onClick(picker) {
            fillDate(picker, new Date() - new Date(new Date().getFullYear().toString()));
        },
    },
    {
        text: "本月",
        onClick(picker) {
            fillDate(picker, (new Date().getDate() - 1) * 24 * 3600 * 1000)
        },
    },
    // {
    //     text: "下月",
    //     onClick(picker) {
    //         fillRelativeDate(picker);
    //     },
    // },
    {
        text: "最近一周",
        onClick(picker) {
            fillDate(picker, passWeek);
        },
    },
    {
        text: "最近一个月",
        onClick(picker) {
            fillDate(picker, passMonth);
        },
    },
    {
        text: "最近三个月",
        onClick(picker) {
            fillDate(picker, passTriMonth);
        },
    },
    {
        text: "最近半年",
        onClick(picker) {
            fillDate(picker, passHalfYear);
        },
    },
    {
        text: "最近一年",
        onClick(picker) {
            fillDate(picker, passYear);
        },
    },
    {
        text: "最近三年",
        onClick(picker) {
            fillDate(picker, passTriYear);
        },
    },
    {
        text: "最近五年",
        onClick(picker) {
            fillDate(picker, passFiveYear);
        },
    },
]

const yearPeriod = 'year'
const monthPeriod = 'month'
const weekPeriod = 'week'
const dayPeriod = 'day'

const getFormatByPeriod = function (period) {
    switch (period) {
        case yearPeriod:
            return "YYYY"
        case monthPeriod:
            return "YYYY-MM"
        case weekPeriod:
            return "YYYY-MM-dd"
        case dayPeriod:
            return "YYYY-MM-dd"
    }
}

export {
    formatter, dateShortCut, getFormatByPeriod, yearPeriod, monthPeriod, dayPeriod, weekPeriod,
    calDate, passWeek, passMonth, passTriMonth, passHalfYear, passYear
}