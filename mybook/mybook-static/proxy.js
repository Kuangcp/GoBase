/* eslint-disable indent */
module.exports = (() => {
    const dev = process.env.npm_config_env
    const protocol = 'http://'
    if (!dev) {
        console.log(`代理的地址是：${process.env.baseUrl}`)
        return `${process.env.baseUrl}`
    }
    const objProxy = {
        lo: `localhost`
    }
    console.log(`代理的地址是：${protocol}${objProxy[dev]}`)
    return `${protocol}${objProxy[dev]}`
})()
