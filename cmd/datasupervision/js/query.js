function executeSQLQuery() {
    var query = document.getElementById("sql-query").value;

    console.log("SQL Query:", query);

    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            console.log(xhr.responseText);
            data = JSON.parse(xhr.responseText);
            console.log('Success:', data);

            const tableInfoDiv = document.createElement('div');
            tableInfoDiv.className = 'table-info';
            tableInfoDiv.id = 'table-info';

            const sqlInputDiv = document.createElement('div');
            sqlInputDiv.className = 'sql-input';

            const sqlQuery = document.getElementById("sql-query").value;

            const textarea = document.createElement('textarea');
            textarea.id = 'sql-query';
            textarea.placeholder = 'Enter SQL query';
            textarea.textContent = sqlQuery;

            const executeButton = document.createElement('button');
            executeButton.textContent = 'Execute';
            executeButton.onclick = executeSQLQuery;

            sqlInputDiv.appendChild(textarea);
            sqlInputDiv.appendChild(executeButton);

            tableInfoDiv.appendChild(sqlInputDiv);

            if (data.Columns && data.Rows) {
                const tableData = buildTableData(data);
                tableInfoDiv.appendChild(tableData);
            } else {
                const noDataMessage = document.createElement('p');
                noDataMessage.textContent = 'No data available to show.';
                tableInfoDiv.appendChild(noDataMessage);
            }

            const tableInfoContainer = document.getElementById('table-info');
            tableInfoContainer.innerHTML = '';
            tableInfoContainer.appendChild(tableInfoDiv);
        } else {
            log(xhr.response);
        }
    };

    var postData = { "query": query };
    var url = `/${schemaName}/sql-input`;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(JSON.stringify(postData));
}

function createSQLInput() {
    const tableInfoDiv = document.createElement('div');
    tableInfoDiv.className = 'table-info';
    tableInfoDiv.id = 'table-info';

    const sqlInputDiv = document.createElement('div');
    sqlInputDiv.className = 'sql-input';

    const textarea = document.createElement('textarea');
    textarea.id = 'sql-query';
    textarea.placeholder = 'Enter SQL query';

    const executeButton = document.createElement('button');
    executeButton.textContent = 'Execute';
    executeButton.onclick = executeSQLQuery;

    sqlInputDiv.appendChild(textarea);
    sqlInputDiv.appendChild(executeButton);
    tableInfoDiv.appendChild(sqlInputDiv);

    document.getElementById('table-info').innerHTML = '';
    document.getElementById('table-info').appendChild(tableInfoDiv);
}