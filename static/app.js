document.addEventListener('DOMContentLoaded', function () {
    loadItems();
    document.getElementById('addItemForm').addEventListener('submit', addItem);
});

function loadItems() {
    fetch('/items')
        .then(response => response.json())
        .then(data => {
            const tableBody = document.getElementById('inventoryTableBody');
            tableBody.innerHTML = '';

            data.forEach(item => {
                const row = document.createElement('tr');

                row.innerHTML = `
                    <td>${item.itemName}</td>
                    <td>
                        <button class="decrement" onclick="updateItem('${item.itemName}', '-')">−</button>
                        ${item.itemQTY}
                        <button class="increment" onclick="updateItem('${item.itemName}', '+')">+</button>
                    </td>
                    <td>${item.itemUsedToDate}</td>
                    <td>${item.minimumQTY}</td>
                    <td>${item.itemType}</td>
                `;

                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error loading items:', error));
}

function addItem(event) {
    event.preventDefault();

    const itemName = document.getElementById('itemName').value.trim();
    const itemType = document.getElementById('itemType').value.trim();
    const itemQTY = document.getElementById('itemQTY').value.trim();
    const minimumQTY = document.getElementById('minimumQTY').value.trim();

    if (!itemName || !itemType || !itemQTY || !minimumQTY) {
        console.error('All fields are required.');
        return;
    }

    const formData = new URLSearchParams();
    formData.append('itemName', itemName);
    formData.append('itemType', itemType);
    formData.append('itemQTY', itemQTY);
    formData.append('minimumQTY', minimumQTY);

    fetch('/item/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
    })
    .then(response => {
        if (response.ok) {
            console.log('Item added successfully.');
            document.getElementById('addItemForm').reset();
            loadItems();
        } else {
            console.error('Failed to add item.');
        }
    })
    .catch(error => console.error('Error adding item:', error));
}

function updateItem(itemName, action) {
    const formData = new URLSearchParams();
    formData.append('itemName', itemName);
    formData.append('action', action);

    fetch('/item/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
    })
    .then(response => {
        if (response.ok) {
            loadItems();
        } else {
            console.error('Failed to update item.');
        }
    })
    .catch(error => console.error('Error updating item:', error));
}
