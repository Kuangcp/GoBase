function loadRecordTables() {
    let typeId = $("#typeId option:selected").val();
    let accountType = $("#accountTypeList option:selected").val();

    let dateQuery = buildWithDefaultDate('startDate', 'endDate', 7);

    url = '/record/list?' + dateQuery + '&typeId=' + typeId + '&accountId=' + accountType;
    handleGet(url, function (data) {
        if (!data.Success) {
            layer.msg('加载账单失败');
            console.log(data);
            return;
        }

        console.log('/category/typeList', data);
        $("#record_table_body tbody").find('tr').each(function () {
            $(this).remove();
        });

        appendRecordRow(data, 'record_table_body');
    });
}

function appendRecordRow(data, targetBlock) {
    let total = 0;
    for (i in data.Data) {
        let record = data.Data[i];

        let line = "<tr>";
        line += '<td>' + record.ID + '</td>';
        line += '<td style="text-align: right"> ' + record.AccountName + '</td>';
        line += '<td>' + record.RecordTypeName + '</td>';
        line += '<td>' + record.CategoryName + '</td>';
        line += '<td style="text-align: right">' + record.Amount / 100.0 + ' </td>';
        line += '<td>' + record.Comment + '</td>';
        line += '<td style="width: 140px">' + record.RecordTime + '</td>';

        line += '</tr>';
        $('#' + targetBlock + ' > tbody:last-child').append(line);
        total += record.Amount;
    }

    $("#total").html('￥' + total / 100.0)
}