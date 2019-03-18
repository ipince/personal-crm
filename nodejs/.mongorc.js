function insertData(dbName, colName, numRecords) {
  var col = db.getSiblingDB(dbName).getCollection(colName);
  for (var i = 0; i < numRecords; i++) {
    col.insert({x: i});
  }
}
