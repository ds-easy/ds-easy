/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "zgdy41tc0kvhqzs",
    "created": "2024-10-18 18:52:35.078Z",
    "updated": "2024-10-18 18:52:35.078Z",
    "name": "Exo_files",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "jxgocgvo",
        "name": "files",
        "type": "file",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "mimeTypes": [],
          "thumbs": [],
          "maxSelect": 1,
          "maxSize": 5242880,
          "protected": false
        }
      }
    ],
    "indexes": [],
    "listRule": null,
    "viewRule": null,
    "createRule": null,
    "updateRule": null,
    "deleteRule": null,
    "options": {}
  });

  return Dao(db).saveCollection(collection);
}, (db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("zgdy41tc0kvhqzs");

  return dao.deleteCollection(collection);
})
