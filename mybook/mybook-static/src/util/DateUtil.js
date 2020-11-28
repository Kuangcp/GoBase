/**
 * 仿照 moment 的接口，简单实现 formatDate
 * @param   {Number} timeStamp    待格式化的 unix 时间戳
 * @return  {String}              格式化后的日期字符串
 */
export default function (timeStamp) {
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
    };
  }