function loadExistConf() {
    get('/queryConfig', function (rspT) {
        let rsp = JSON.parse(rspT)
        if (rsp.code !== 0) {
            alert('查询配置失败：' + rsp.msg)
            return
        }
        let data = rsp.data;

        document.getElementById('id').value = data.id

        // proxy and direct
        document.getElementById('proxy').innerHTML = ''
        if (data.proxy !== null) {
            document.getElementById('proxy-switch-type').checked = data.proxy.proxy_type === 1
            if (data.proxy.paths !== null) {
                for (u of data.proxy.paths) {
                    let tmp = document.createElement('input');
                    tmp.name = 'paths';
                    tmp.value = u;
                    document.getElementById('proxy').appendChild(tmp);
                }
            }
        }

        document.getElementById('direct').innerHTML = ''
        if (data.direct !== null) {
            document.getElementById('direct-switch-type').checked = data.direct.proxy_type === 1
            if (data.direct.paths !== null) {
                for (u of data.direct.paths) {
                    let tmp = document.createElement('input');
                    tmp.name = 'paths';
                    tmp.value = u;
                    document.getElementById('direct').appendChild(tmp);
                }
            }
        }

        // redis
        if (data.redis !== null) {
            for (let child of document.getElementById('redis').children) {
                if (child.tagName === 'INPUT') {
                    child.value = data.redis[child.name]
                }
            }
        }

        // groups
        document.getElementById('groups').innerHTML = ''
        if (data.groups !== null) {
            for (let group of data.groups) {
                addGroup(group);
            }
        }
    })
}
