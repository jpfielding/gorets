import moment from 'moment';
import LZString from 'lz-string';

const cacheVersion = 'v2';

class StorageCache {
  constructor(storage) {
    this.storage = storage;
  }

  clearAll() {
    this.storage.clear();
  }

  versioned(key) {
    return `${key}.${cacheVersion}`;
  }

  remove(k) {
    const key = this.versioned(k);
    this.storage.removeItem(this.versioned(key));
  }

  readCache(k) {
    const key = this.versioned(k);
    const storedRaw = this.storage.getItem(key);
    const decompressedRaw = LZString.decompressFromUTF16(storedRaw);
    if (decompressedRaw) {
      return JSON.parse(decompressedRaw);
    }
    return JSON.parse(storedRaw);
  }

  getFromCache(k) {
    const key = this.versioned(k);
    const storedItem = this.readCache(key);
    if (storedItem == null) {
      console.log('item not found in cache', key);
      return null;
    }
    if (moment().isAfter(storedItem.expiresAt)) {
      this.storage.removeItem(key);
      console.log('item found but expired in cache', key);
      return null;
    }
    console.log('item found in cache', key);
    return storedItem.value;
  }


  putInCache(k, item, ttlMinutes) {
    if (item == null) {
      console.log('rejected null cache item', k, item);
      return null;
    }
    const key = this.versioned(k);
    this.storage.setItem(key, LZString.compressToUTF16(JSON.stringify({
      expiresAt: moment().add(ttlMinutes, 'minutes'),
      value: item,
    })));
    console.log('item stored in cache', key, item);
    this.printCache();
    return item;
  }

  printCache() {
    for (let i = 0; i < this.storage.length; i++) {
      const key = this.storage.key(i);
      const storageItem = this.readCache(key);
      console.log('cache entry:', key, storageItem);
    }
  }
  // clearExpired() {
  //   const clear = () => {
  //     const keysToRemove = [];
  //     for (let i = 0; i < this.storage.length; i++) {
  //       const key = this.storage.key(i);
  //       const storageItem = this.readCache(key);
  //       if (moment().isAfter(storageItem.expiresAt)) {
  //         keysToRemove.push(key);
  //       }
  //     }
  //     keysToRemove.forEach(k => this.storage.removeItem(k));
  //   };
  //   setTimeout(clear, 0);
  // }
}

export default new StorageCache(sessionStorage);
