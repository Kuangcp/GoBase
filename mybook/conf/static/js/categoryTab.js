function loadCategoryTables() {
    let typeId = $("#typeIdMonth option:selected").val();

    url = '/record/month?' + buildWithDefaultDate('startDateMonth', 'endDateMonth', 7) + '&typeId=' + typeId;
    handleGet(url, function (data) {
        if (!data.Success) {
            layer.msg('加载账单失败');
            console.log(data);
            return;
        }

        console.log('/record/month', data);
        $("#month_table_body tbody").find('tr').each(function () {
            $(this).remove();
        });
        appendCategoryRecord(data);
    });
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
        line += '<td style="text-align: right;width: 30px;">' + record.Amount / 100.0 + ' </td>';
        line += '<td style="text-align: right;width: 120px;">' + record.Date + '</td>';
        line += '<td style="width: 50px;"> <button onclick="loadCategoryRecordDetail(' + record.CategoryId + ')">详情</button></td>';
        line += '<td style="width: 80px;"> <button onclick="loadCategoryRecordWeek(' + record.CategoryId + ')">周统计</button></td>';

        line += '</tr>';
        $('#month_table_body > tbody:last-child').append(line);
        total += record.Amount;
    }
}

// 分类详情数据
function loadCategoryRecordDetail(category) {
    tip(['750px', '420px'], '单分类明细账单', $("#month_detail_tables").html());

    let dateQuery = buildWithDefaultDate('startDateMonth', 'endDateMonth', 7);
    handleGet('/record/monthDetail?' + dateQuery + '&categoryId=' + category, function (data) {
        if (data.Success) {
            console.log('/record/monthDetail', data);

            appendRecordRow(data, 'month_detail_table_body')
        } else {
            layer.msg('加载分类明细失败');
            console.log(data)
        }
    });
}

// 分类详情周统计数据
function loadCategoryRecordWeek(category) {
    tip(['750px', '420px'], '分类周明细账单', $("#month_detail_tables").html());
    let dateQuery = buildWithDefaultDate('startDateMonth', 'endDateMonth', 7);
    handleGet('/record/monthDetail?' + dateQuery + '&categoryId=' + category, function (data) {
        if (data.Success) {
            console.log('/record/monthDetail', data);

            appendRecordRow(data, 'month_detail_table_body')
        } else {
            layer.msg('加载分类明细失败');
            console.log(data)
        }
    });
}
