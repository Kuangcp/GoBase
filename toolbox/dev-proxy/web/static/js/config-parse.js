function addRoute(here) {
    let tmp = document.createElement('div');
    tmp.className = 'routers'
    tmp.innerHTML = document.getElementById('router-temp').innerHTML
    here.parentNode.parentElement.appendChild(tmp)
}

function addGroup(val) {
    let tmp = document.createElement('div');
    tmp.className = 'group'
    tmp.innerHTML = document.getElementById('group-temp').innerHTML
    if (val !== null && val !== undefined) {
        let nav = tmp.children[0]
        for (let child of nav.children) {
            if (child.name === 'proxy_type') {
                child.checked = val.proxy_type === 1
            } else if (child.name === 'name') {
                child.value = val.name
            }
        }
        for (let router of val.routers) {
            let routeDom = document.createElement('div');
            routeDom.className = 'routers'
            routeDom.innerHTML = document.getElementById('router-temp').innerHTML
            for (let child of routeDom.children) {
                if (child.name === 'proxy_type') {
                    child.checked = router.proxy_type === 1
                } else if (child.name === 'src') {
                    child.value = router.src
                } else if (child.name === 'dst') {
                    child.value = router.dst
                }
            }
            tmp.appendChild(routeDom)
        }
    }
    document.getElementById('groups').appendChild(tmp)
}

function addProxyPath(divId) {
    let tmp = document.createElement('input');
    tmp.name = 'paths'
    document.getElementById(divId).appendChild(tmp)
}

function save() {
    let config = {}
    config['id'] = document.getElementById('id').value
    config['groups'] = buildGroups()
    config['redis'] = divToObj('redis')
    let proxy = divToObj('proxy');
    proxy.proxy_type = document.getElementById('proxy-switch-type').checked ? 1 : 0
    config['proxy'] = proxy
    let direct = divToObj('direct');
    direct.proxy_type = document.getElementById('direct-switch-type').checked ? 1 : 0
    config['direct'] = direct

    if (config.redis.db === '') {
        config.redis.db = null
    }
    if (config.redis.pool_size === '') {
        config.redis.pool_size = null
    }
    config.proxy.name = '抓包'
    config.direct.name = '直连'
    postJSON('/saveConfig', JSON.stringify(config), function (data) {
        let rsp = JSON.parse(data);
        if (rsp.code !== 0) {
            alert(rsp.msg)
        }
    })
}

function buildGroups() {
    let groups = []

    let groupNodes = document.getElementById('groups').children;
    for (let i = 0; i < groupNodes.length; i++) {
        let groupNode = groupNodes[i];

        let group = buildGroup(groupNode);
        groups.push(group)
    }

    return groups
}

function buildGroup(node) {
    let fields = []
    let routers = []
    for (let tmp of node.children) {
        if (tmp.className === 'group-nav') {
            for (let child of tmp.children) {
                if (child.tagName === 'INPUT') {
                    fields.push(child)
                }
            }
        } else if (tmp.className === 'routers') {
            let route = inputToObj(tmp.getElementsByTagName('input'));
            routers.push(route)
        }
    }

    let group = inputToObj(fields);
    group['routers'] = routers
    return group
}

function divToObj(id) {
    return inputToObj(document.getElementById(id).getElementsByTagName('input'))
}

function toInt(val) {
    return val * 1
}

function inputToObj(fields) {
    let result = {}
    for (let i = 0; i < fields.length; i++) {
        let field = fields[i];
        let name = field.name;
        if (!(name in result)) {
            if (name === 'proxy_type') {
                result[name] = field.checked ? 1 : 0
            } else {
                if (field.type === 'number') {
                    result[name] = toInt(field.value);
                } else {
                    result[name] = field.value;
                }
            }
        } else {
            let last = result[name]
            if (last instanceof Array) {
                last.push(field.value)
            } else {
                let tmp = []
                tmp.push(last, field.value)
                result[name] = tmp
            }
        }
    }
    return result
}