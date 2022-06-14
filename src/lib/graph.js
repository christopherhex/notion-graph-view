
import { local } from "d3";
import wasmLib from "./wasm";

/**
 * 
 * @param {String} notionKey - add notion Key
 * @param {Boolean} forceRefresh - force a refresh of data instead of
 */
const getGraphData = async (notionKey, forceRefresh = false) => {

    const LOCALSTORAGE_DATA_KEY = 'notionGraphData';

    const storageAvailable = typeof (Storage) !== "undefined";

    // Check if data can be found in localstorage
    if (storageAvailable && !forceRefresh) {

        // Try to find in local storage
        const localData = localStorage.getItem(LOCALSTORAGE_DATA_KEY);

        if (localData) {
            return JSON.parse(localData);
        }

    }

    // Make sure Wasm is initiated
    await wasmLib();
    const newData = await window.getData(notionKey);

    if (storageAvailable) {
        localStorage.setItem(LOCALSTORAGE_DATA_KEY, newData)
    }
    console.log("NEW DATA")
    console.log(newData);
    return JSON.parse(newData);

}


export default getGraphData;