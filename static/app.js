document.addEventListener('DOMContentLoaded', function () {
    loadItems();
    document.getElementById('addItemForm').addEventListener('submit', addItem);
});

/**
 * loadItems fetches the list of inventory items and populates the table.
 */
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
                    <td>${item.itemTotalTossed || 0}</td> <!-- Total Tossed -->
                    <td>${item.minimumQTY}</td>
                    <td>${item.itemTypeName}</td>
                    <td>${item.itemSubstitutionName}</td>
                    <td>
                        <button class="dispose" onclick="disposeItem('${item.itemName}')">🗑️ Dispose</button>
                    </td>
                `;

                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error loading items:', error));
}

/**
 * addItem handles form submission to add a new inventory item.
 */
function addItem(event) {
    event.preventDefault();

    const itemName = document.getElementById('itemName').value.trim();
    const itemTypeID = document.getElementById('itemTypeID').value;
    const itemSubstitutionID = document.getElementById('itemSubstitutionID').value;
    const itemQTY = document.getElementById('itemQTY').value.trim();
    const minimumQTY = document.getElementById('minimumQTY').value.trim();
    const itemExpirationPeriod = document.getElementById('itemExpirationPeriod').value.trim();

    if (!itemName || !itemTypeID || !itemSubstitutionID || !itemQTY || !minimumQTY || !itemExpirationPeriod) {
        console.error('All fields are required.');
        return;
    }

    const formData = new URLSearchParams();
    formData.append('itemName', itemName);
    formData.append('itemTypeID', itemTypeID);
    formData.append('itemSubstitutionID', itemSubstitutionID);
    formData.append('itemQTY', itemQTY);
    formData.append('minimumQTY', minimumQTY);
    formData.append('itemExpirationPeriod', itemExpirationPeriod);

    fetch('/item/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('addItemForm').reset();
            loadItems();
        } else {
            console.error('Failed to add item.');
        }
    })
    .catch(error => console.error('Error adding item:', error));
}

/**
 * updateItem sends a request to update the quantity of an inventory item.
 */
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
            console.error('Failed to update item.', response.statusText);
        }
    })
    .catch(error => console.error('Error updating item:', error));
}

/**
 * disposeItem sends a request to dispose of an expired inventory item.
 */
function disposeItem(itemName) {
    const formData = new URLSearchParams();
    formData.append('itemName', itemName);

    fetch('/item/dispose', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: formData.toString()
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Failed to dispose item.');
        }
    })
    .then(data => {
        console.log('Disposed:', data);
        loadItems(); // Refresh the table, no alert
    })
    .catch(error => console.error('Error disposing item:', error));
}
