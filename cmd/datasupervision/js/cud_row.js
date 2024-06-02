function insertRow(tableName, columnsString) {
    var columns = JSON.parse(decodeURIComponent(columnsString));
    console.log("insertRow: columns:", columns);

    var row = new Map();
    columns.forEach(function (column, index) {
        var insertValue = document.getElementById("insert-value-" + index).value;
        row.set(column, insertValue);
    });

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log(xhr.responseText);
            data = JSON.parse(xhr.responseText);
            console.log('Success:', data);

            getTableData(tableName);
        } else {
            log(xhr.response);
        }
    };

    var rowObject = Object.fromEntries(row);

    var postData = JSON.stringify({ "row": rowObject });
    console.log("postData", postData);

    var url = `/${schemaName}/insert-row?tableName=` + tableName;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(postData);
}


function editRow(tableName, rowIndex, tableDataString) {
    var tableData = JSON.parse(decodeURIComponent(tableDataString));
    var editTable = buildTableData(tableData, tableName, rowIndex, tableDataString);

    const tableInfoContainer = document.getElementById("table-content");
    tableInfoContainer.innerHTML = '';
    tableInfoContainer.appendChild(editTable);
}

function updateRow(tableName, rowIndex, tableDataString, notUpdate) {
    var tableData = JSON.parse(decodeURIComponent(tableDataString));
    console.log("updateRow", tableName, rowIndex, tableData)

    if (notUpdate) {
        getTableData(tableName);
        return
    }

    var row = {
        "columns": [],
        "oldRow": [],
        "newRow": []
    };
    tableData.Columns.forEach(function (column, index) {
        var newValue = document.getElementById("update-value-" + index).value;
        var oldValue = tableData.Rows[rowIndex][index];
        console.log("updateValue", index, oldValue, newValue);
        row.columns.push(column);
        row.oldRow.push(oldValue);
        row.newRow.push(newValue);
    });

    console.log("row", row);

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            data = JSON.parse(xhr.responseText);
            console.log('Success:', data);

            getTableData(tableName);
        } else {
            log(xhr.response);
        }
    };

    var postData = JSON.stringify(row);
    console.log("postData", postData);

    var url = `/${schemaName}/update-row?tableName=` + tableName;
    xhr.open("PATCH", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(postData);
}

function deleteRow(tableName, rowIndex, tableDataString) {
    var tableData = JSON.parse(decodeURIComponent(tableDataString));
    console.log("deleteRow", tableName, rowIndex, tableData)

    var row = new Map();
    tableData.Columns.forEach(function (column, index) {
        var deleteValue = tableData.Rows[rowIndex][index];
        console.log("insertValue", index, deleteValue);
        row.set(column, deleteValue);
    });

    console.log("row", row);

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            data = JSON.parse(xhr.responseText);
            console.log('Success:', data);

            getTableData(tableName);
        } else {
            log(xhr.response);
        }
    };

    var rowObject = Object.fromEntries(row);
    var postData = JSON.stringify({ "row": rowObject });
    console.log("postData", postData);

    var url = `/${schemaName}/delete-row?tableName=` + tableName;
    xhr.open("PATCH", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(postData);
}