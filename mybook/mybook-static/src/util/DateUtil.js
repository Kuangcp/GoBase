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
            let now = new Date();
            let passDay = (now.getDate() - 1) * 24 * 3600 * 1000
            const start = now;
            start.setTime(start.getTime() - passDay)
            picker.$emit("pick", [start, new Date()]);
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
            fillDate(picker, 3600 * 1000 * 24 * 7);
        },
    },
    {
        text: "最近一个月",
        onClick(picker) {
            fillDate(picker, 3600 * 1000 * 24 * 30);
        },
    },
    {
        text: "最近三个月",
        onClick(picker) {
            fillDate(picker, 3600 * 1000 * 24 * 90);
        },
    },
    {
        text: "最近半年",
        onClick(picker) {
            fillDate(picker, 3600 * 1000 * 24 * 180);
        },
    },
    {
        text: "最近一年",
        onClick(picker) {
            fillDate(picker, 3600 * 1000 * 24 * 365);
        },
    },
    {
        text: "最近三年",
        onClick(picker) {
            fillDate(picker, 3600 * 1000 * 24 * 365 * 3);
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
    formatter, dateShortCut, getFormatByPeriod, yearPeriod, monthPeriod, dayPeriod, weekPeriod
}