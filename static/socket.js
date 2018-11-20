'use strict';
const conn = new WebSocket("ws://" + document.location.host + "/ws");

// Need to implement some Redux here in this thing...
(() => {
    const table = document.querySelector('#mytable');
    const tbody = document.querySelector('tbody');
    const insertButton = table.querySelector('[action="insert"]');
    const tmplRow = document.querySelector('[id="user-row"]');
    const tmplCell = document.querySelector('[id="cell"]');

    const processCell = (row, key, data) => {
        const td = row.querySelector(`[key="${key}"]`);
        const cell = document.importNode(tmplCell.content, true);
        const span = cell.querySelector('span');
        const input = cell.querySelector('input');
        const form = cell.querySelector('form')

        span.textContent = data[key] ? data[key] : '-';
        input.value = data[key] ? data[key] : '';
        input.name = key;
        cell.querySelector('input[name="ID"]').value = data.ID;

        const toggleSpanToForm = () => {
            span.hidden = true;
            form.hidden = false;
        };

        span.addEventListener('click', toggleSpanToForm);

        form.addEventListener('submit', e => {
            e.preventDefault();
            span.hidden = false;
            form.hidden = true;

            td.setAttribute('disabled', '');
            span.removeEventListener('click', toggleSpanToForm);

            let fd;
            const data = new FormData(e.target).toJSON();
            if (data.ID !== "undefined") {
                data.ID = Number(data.ID);
                span.textContent = data[key] ? data[key] : '-';

                fd = {
                    Jsonrpc: "2.0",
                    Method: 'updateUser',
                    Params: data
                };
            } else {
                //data.ID = null;
                fd = {
                    Jsonrpc: "2.0",
                    Method: 'insertUser',
                    Params: data
                };
            }
            conn.send(JSON.stringify(fd));
        });

        td.appendChild(cell);
    };

    const processRow = (data, row = document.importNode(tmplRow.content, true)) => {
        row.querySelector('tr').setAttribute('ID', data.ID)
        processCell(row, 'Name', data);
        processCell(row, 'Email', data);
        return row;
    }

    insertButton.insertingNew = false;
    insertButton.addEventListener('click', () => {
        if (!insertButton.insertingNew) {
            tbody.appendChild(processRow({}));
        }
        insertButton.insertingNew = true;
    });

    conn.onmessage = (e) => {
        const data = JSON.parse(e.data);
        const result = data.Result;
        if (data.ID === 'id-for-getUsers') {
            result.forEach(element => tbody.appendChild(processRow(element)));
        } else {
            let row = tbody.querySelector(`tr[id="${result.ID}"]`);
            if (!row) {
                row = tbody.querySelector(`tr[id="undefined"]`);
                if (!row) {
                    return tbody.appendChild(processRow(result));
                }
                row.setAttribute('ID', result.ID)
                insertButton.insertingNew = false;
            }
            Object.keys(result).filter(key => key != 'ID').forEach(key => {
                const cell = row.querySelector(`[key="${key}"]`);
                cell.innerHTML = '';
                cell.removeAttribute('disabled');
                processCell(row, key, result);
            });
        }
    }
    conn.onopen = () => {
        const message = {
            Jsonrpc: '2.0',
            Method: 'getUsers',
            ID: 'id-for-getUsers'
        };
        conn.send(JSON.stringify(message));
    }
})();