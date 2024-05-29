document.addEventListener('DOMContentLoaded', function() {
    const dbTypeSelector = document.getElementById('dbType');
    const hostInput = document.getElementById('host');
    const portInput = document.getElementById('port');
    const databaseInput = document.getElementById('database');
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');

    const config = {
        PostgreSQL: {
            host: 'localhost',
            port: '5432',
            database: 'blog',
            username: 'postgres',
            password: 'postgres'
        },
        MySQL: {
            host: 'localhost',
            port: '3306',
            database: 'mysql',
            username: 'root',
            password: 'password'
        },
    };

    function updateFields() {
        const selectedDb = dbTypeSelector.value;
        const settings = config[selectedDb];

        hostInput.value = settings.host;
        portInput.value = settings.port;
        databaseInput.value = settings.database;
        usernameInput.value = settings.username;
        passwordInput.value = settings.password;
    }

    dbTypeSelector.addEventListener('change', updateFields);
    updateFields();
});