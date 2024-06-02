document.addEventListener('DOMContentLoaded', function () {
    const openPopupButtons = document.querySelectorAll('.openPopupBtn');
    const popups = document.querySelectorAll('.popup');
    const closePopupButtons = document.querySelectorAll('.close');

    openPopupButtons.forEach((btn, index) => {
        btn.addEventListener('click', function () {
            const popup = document.getElementById(`popup-${index}`);
            popup.style.display = 'flex';
        });
    });

    closePopupButtons.forEach((btn, index) => {
        btn.addEventListener('click', function () {
            const popup = document.getElementById(`popup-${index}`);
            popup.style.display = 'none';
        });
    });

    window.addEventListener('click', function (event) {
        popups.forEach((popup, index) => {
            if (event.target === popup) {
                popup.style.display = 'none';
            }
        });
    });
});

function downloadTableDataJson(tableName) {
    fetch(`/${schemaName}/download/json/` + tableName)
        .then(response => response.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = 'data.json';
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        })
        .catch(err => console.error('Download error:', err));
}

function handleFileImportJson(tableName) {
    console.log("handleFileImportJson");
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.json';
    fileInput.onchange = () => {
        const file = fileInput.files[0];
        if (file) {
            uploadFile(file);
        }
    };
    fileInput.click();
}

function uploadFile(file) {
    console.log("uploadFile");
    const formData = new FormData();
    formData.append('file', file);

    fetch(`/${schemaName}/insert-rows/json`, {
        method: 'POST',
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            console.log('File uploaded successfully', data);
        })
        .catch(error => {
            console.error('Error uploading file:', error);
        });
}

function downloadTableDataCSV(tableName) {
    console.log("test 2 downloadTableDataCVS");
    fetch(`/${schemaName}/download/csv/` + tableName)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.blob();
        })
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.style.display = 'none';
            a.href = url;
            a.download = tableName + '.csv';
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
        })
        .catch(error => console.error('There was a problem with the fetch operation:', error));
}

function handleFileImportCSV(tableName) {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv';
    fileInput.onchange = () => {
        const file = fileInput.files[0];
        if (file) {
            uploadFileCSV(file, tableName);
        }
    };
    fileInput.click();
}

function uploadFileCSV(file, tableName) {
    const formData = new FormData();
    formData.append('file', file);

    fetch(`/${schemaName}/insert-rows/csv?tableName=${tableName}`, {
        method: 'POST',
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            console.log('File uploaded successfully', data);
        })
        .catch(error => {
            console.error('Error uploading file:', error);
        });
}
