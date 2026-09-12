<?php
// api.php
// A simple local SQLite backend to keep data private and off the browser's localStorage.

header("Content-Type: application/json");
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Methods: GET, POST, OPTIONS");
header("Access-Control-Allow-Headers: Content-Type");

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    exit(0);
}

// Ensure the database file is saved in the same directory and is excluded via .gitignore
$db_file = __DIR__ . '/private_local_data.sqlite';
$db = new SQLite3($db_file);

// Create a simple key-value store table
$db->exec("CREATE TABLE IF NOT EXISTS store (key TEXT PRIMARY KEY, value TEXT)");

$method = $_SERVER['REQUEST_METHOD'];

if ($method === 'GET') {
    $key = $_GET['key'] ?? '';
    if (!$key) { 
        echo json_encode(["error" => "No key provided"]); 
        exit; 
    }
    
    $stmt = $db->prepare("SELECT value FROM store WHERE key = :key");
    $stmt->bindValue(':key', $key, SQLITE3_TEXT);
    $res = $stmt->execute()->fetchArray(SQLITE3_ASSOC);
    
    if ($res) {
        echo $res['value']; // Value is already a JSON string
    } else {
        echo "null";
    }
} 
elseif ($method === 'POST') {
    $input = json_decode(file_get_contents("php://input"), true);
    $key = $input['key'] ?? '';
    $value = $input['value'] ?? '';
    
    if (!$key) { 
        echo json_encode(["error" => "No key provided"]); 
        exit; 
    }
    
    // Convert array/object back to string if needed
    $valStr = is_string($value) ? $value : json_encode($value);

    // REPLACE INTO works safely across all SQLite versions
    $stmt = $db->prepare("REPLACE INTO store (key, value) VALUES (:key, :value)");
    $stmt->bindValue(':key', $key, SQLITE3_TEXT);
    $stmt->bindValue(':value', $valStr, SQLITE3_TEXT);
    $stmt->execute();
    
    echo json_encode(["status" => "success"]);
}
?>
