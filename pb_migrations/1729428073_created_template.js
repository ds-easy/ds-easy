/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const collection = new Collection({
    "id": "2pw7qlnaim53zrj",
    "created": "2024-10-20 12:41:13.404Z",
    "updated": "2024-10-20 12:41:13.404Z",
    "name": "template",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "ottxtsat",
        "name": "file",
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
  const collection = dao.findCollectionByNameOrId("2pw7qlnaim53zrj");

  return dao.deleteCollection(collection);
})
