import moment from 'moment';
import LZString from 'lz-string';

const cacheVersion = 'v2';

const getKey = () => `rets-storage.${cacheVersion}`;

class StorageCache {
  constructor(storage) {
    this.storage = storage;
  }

  clearAll() {
    this.storage.clear();
  }

  remove() {
    this.storage.removeItem(getKey());
  }

  // cachedGet(supplier, ttl) {
  //   const key = getKey();
  //   const cachedValue = this.getFromCache(key);
  //   if (cachedValue) {
  //     return cachedValue;
  //   }
  //
  //   const value = supplier.apply();
  //   if (value) {
  //     setTimeout(() => this.putInCache(value, ttl, key), 0);
  //   }
  //   return value;
  // }

  readCache(key) {
    const storedRaw = this.storage.getItem(key);
    const decompressedRaw = LZString.decompressFromUTF16(storedRaw);
    if (decompressedRaw) {
      return JSON.parse(decompressedRaw);
    }
    return JSON.parse(storedRaw);
  }

  getFromCache() {
    const key = getKey();

    const storedItem = this.readCache(key);
    if (storedItem == null) {
      return null;
    }

    if (moment().isAfter(storedItem.expiresAt)) {
      this.storage.removeItem(key);
      return null;
    }

    return storedItem.value;
  }


  putInCache(item, ttlMinutes) {
    const key = getKey();

    this.storage.setItem(key, LZString.compressToUTF16(JSON.stringify({
      expiresAt: moment().add(ttlMinutes, 'minutes'),
      value: item,
    })));

    return item;
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
