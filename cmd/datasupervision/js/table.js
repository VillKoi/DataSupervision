const schemaName = window.config.schemaName;

function startTransaction() {
    console.log("startTransaction");

    document.getElementById("btn-transaction").style.display = 'none';
    document.getElementById("btn-rollback").style.display = 'inline';
    document.getElementById("btn-commit").style.display = 'inline';

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var tableData = JSON.parse(xhr.responseText);
            log("success")
        } else {
            log(xhr.response);
        }
    };

    xhr.open("GET", "/${schemaName}/begin-transaction", true);
    xhr.send();
}

function Rollback() {
    console.log("Rollback");

    document.getElementById("btn-transaction").style.display = 'inline';
    document.getElementById("btn-rollback").style.display = 'none';
    document.getElementById("btn-commit").style.display = 'none';

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var tableData = JSON.parse(xhr.responseText);
            log("success")
        } else {
            log(xhr.response);
        }
    };

    xhr.open("GET", "/${schemaName}/rollback", true);
    xhr.send();
}

function Commit() {
    console.log("Commit");

    document.getElementById("btn-transaction").style.display = 'inline';
    document.getElementById("btn-rollback").style.display = 'none';
    document.getElementById("btn-commit").style.display = 'none';

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var tableData = JSON.parse(xhr.responseText);
            log("success")
        } else {
            log(xhr.response);
        }
    };

    xhr.open("GET", "/${schemaName}/commit", true);
    xhr.send();
}



function getTableData(tableName) {
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            var tableData = JSON.parse(xhr.responseText);
            displayTableData(tableName, tableData);
        }
    };
    xhr.open("GET", "/${schemaName}/tabledata/" + tableName, true);
    xhr.send();
}