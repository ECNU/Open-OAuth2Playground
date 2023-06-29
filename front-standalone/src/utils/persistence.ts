interface StorageItem {
    key: string;
    value: any;
}

export class LocalStorageService {
    private storage: Storage;

    constructor() {
        this.storage = window.localStorage;
    }

    addItem(item: StorageItem): void {
        const existingItem = this.getItem(item.key);
        if (existingItem) {
            this.updateItem(item);
        } else {
            this.storage.setItem(item.key, JSON.stringify(item.value));
        }
    }

    getItem(key: string): any {
        const item = this.storage.getItem(key);
        return item ? JSON.parse(item) : null;
    }

    updateItem(item: StorageItem): void {
        const existingItem = this.getItem(item.key);
        if (existingItem) {
            this.storage.setItem(item.key, JSON.stringify(item.value));
        } else {
            this.addItem(item);
        }
    }

    removeItem(key: string): void {
        this.storage.removeItem(key);
    }
}