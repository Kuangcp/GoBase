function loadCategoryTables() {
    let typeId = $("#typeIdMonth option:selected").val();

    url = '/record/category?' + buildWithDefaultDate('startDateMonth', 'endDateMonth', 7) + '&typeId=' + typeId;
    handleGet(url, function (data) {
        if (!data.Success) {
            layer.msg('加载账单失败');
            console.log(data);
            return;
        }

        console.log('/record/category', data);
        $("#category_tables tbody").find('tr').each(function () {
            $(this).remove();
        });
        appendCategoryRecord(data);
    });
}

function buildFormatPrice(amount) {
    first = parseInt(amount / 100)
    secon = amount % 100
    if (secon < 0) {
        secon *= -1
    }
    if (secon === 0) {
        secon = '00';
    } else if (secon < 10) {
        secon += '0';
    }


    return first + '.' + secon
}

// 分类数据
function appendCategoryRecord(data) {
    let total = 0;
    for (i in data.Data) {
        let record = data.Data[i];

        let line = "<tr>";
        line += '<td>' + record.CategoryId + '</td>';
        line += '<td>' + record.RecordTypeName + '</td>';
        line += '<td style="text-align: right;width: 30px;"> ' + record.Name + '</td>';
        line += '<td style="text-align: right;width: 30px;">' + buildFormatPrice(record.Amount) + ' </td>';
        line += '<td style="text-align: right;width: 120px;">' + record.Date + '</td>';
        line += '<td style="width: 50px;"> <button onclick="loadCategoryRecordDetail('
            + record.CategoryId + ', \'' + record.Name + '\')">详情</button></td>';
        line += '<td style="width: 80px;"> <button onclick="loadCategoryRecordWeek('
            + record.CategoryId + ', \'' + record.Name + '\')">周统计</button></td>';
        line += '<td style="width: 80px;"> <button onclick="loadCategoryRecordMonth('
            + record.CategoryId + ', \'' + record.Name + '\')">月统计</button></td>';
        line += '</tr>';
        $('#category_tables > tbody:last-child').append(line);
        total += record.Amount;
    }
    $('#categoryTotal').html('￥' + buildFormatPrice(total))
}

// 分类详情数据
function loadCategoryRecordDetail(category, name) {
    tip(['750px', '420px'], '明细账单 - ' + name, $("#category_detail_tables").html());
    let typeId = $("#typeIdMonth option:selected").val();
    if (typeId === '3' || typeId === "4") {
        typeId = ""
    }

    let dateQuery = buildWithDefaultDate('startDateMonth', 'endDateMonth', 7);
    handleGet('/record/categoryDetail?' + dateQuery + '&categoryId=' + category + "&typeId=" + typeId,
        function (data) {
            if (data.Success) {
                console.log('/record/categoryDetail', data);

                appendRecordRow(data, 'category_detail_table')
            } else {
                layer.msg('加载分类明细失败');
                console.log(data)
            }
        });
}

// 分类详情周统计数据
function loadCategoryRecordWeek(category, name) {
    tip(['320px', '450px'], '周明细账单 - ' + name, $("#category_week_detail_tables").html());
    let dateQuery = buildWithDefaultDate('startDateMonth', 'endDateMonth', 7);
    let typeId = $("#typeIdMonth option:selected").val();
    handleGet('/record/categoryWeekDetail?' + dateQuery + '&categoryId=' + category + "&typeId=" + typeId,
        categoryRecordGroupData());
}

// 分类详情月统计数据
function loadCategoryRecordMonth(category, name) {
    tip(['320px', '450px'], '月明细账单 - ' + name, $("#category_week_detail_tables").html());
    let dateQuery = buildWithDefaultDate('startDateMonth', 'endDateMonth', 7);
    let typeId = $("#typeIdMonth option:selected").val();
    handleGet('/record/categoryMonthDetail?' + dateQuery + '&categoryId=' + category + "&typeId=" + typeId,
        categoryRecordGroupData());
}

function categoryRecordGroupData() {
    return function (data) {
        if (data.Success) {
            console.log('/record/categoryWeekDetail', data);

            let total = 0;
            for (i in data.Data) {
                let record = data.Data[i];

                let line = "<tr>";
                line += '<td style="text-align: center;width: 110px;">' + record.StartDate + '</td>';
                line += '<td style="text-align: center;width: 110px;">' + record.EndDate + '</td>';
                line += '<td style="text-align: right;width: 30px;">' + buildFormatPrice(record.Amount) + '</td>';
                line += '</tr>';

                total += record.Amount;
                $('#category_week_detail_table_body > tbody:last-child').append(line);
            }
        } else {
            layer.msg('加载分类明细失败');
            console.log(data)
        }
    };
}
