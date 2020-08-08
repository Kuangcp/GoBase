function calculateAccountBalance() {
    $("#account_balance_table tbody").find('tr').each(function () {
        $(this).remove();
    });

    handleGet('/account/balance', function (data) {
        if (!data.Success) {
            layer.msg('加载帐户列表失败');
            console.log(data);
            return;
        }

        console.log('/account/balance', data);
        for (i in data.Data) {
            // handleAccount(data.Data[i])
            appendAccountRecordRow(data.Data[i], 'account_balance_table');
        }
    });
}

function appendAccountRecordRow(record, targetBlock) {
    let line = "<tr>";
    line += '<td>' + record.ID + '</td>';
    line += '<td style="text-align: right"> ' + record.Name + '</td>';
    line += '<td style="text-align: right">' + buildFormatPrice(record.InitAmount) + ' </td>';
    line += '<td style="text-align: right">' + buildFormatPrice(record.CurrentAmount) + ' </td>';
    line += '</tr>';
    $('#' + targetBlock + ' > tbody:last-child').append(line);
}