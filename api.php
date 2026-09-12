<?php
// api.php
// A proper relational SQLite backend for CommandSpeak

header("Content-Type: application/json");
header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Methods: GET, POST, OPTIONS");
header("Access-Control-Allow-Headers: Content-Type");

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    exit(0);
}

$db_file = __DIR__ . '/private_local_data.sqlite';
$db = new SQLite3($db_file);

// --- 1. PROPER RELATIONAL SCHEMA ---
$db->exec("
    CREATE TABLE IF NOT EXISTS config (
        key TEXT PRIMARY KEY,
        value TEXT
    );
    CREATE TABLE IF NOT EXISTS repositories (
        id TEXT PRIMARY KEY,
        name TEXT,
        path TEXT,
        active INTEGER
    );
    CREATE TABLE IF NOT EXISTS history (
        id TEXT PRIMARY KEY,
        repoId TEXT,
        query TEXT,
        intent TEXT,
        cmd TEXT,
        timestamp TEXT,
        dryRun INTEGER
    );
");

$method = $_SERVER['REQUEST_METHOD'];

if ($method === 'GET') {
    $key = $_GET['key'] ?? '';
    if (!$key) { echo json_encode(["error" => "No key provided"]); exit; }
    
    if ($key === 'commandspeak_config') {
        $stmt = $db->prepare("SELECT value FROM config WHERE key = 'user_config'");
        $res = $stmt->execute()->fetchArray(SQLITE3_ASSOC);
        echo $res ? $res['value'] : "null";
    }
    elseif ($key === 'commandspeak_repos') {
        $res = $db->query("SELECT * FROM repositories");
        $repos = [];
        while ($row = $res->fetchArray(SQLITE3_ASSOC)) {
            $row['active'] = (bool)$row['active'];
            $repos[] = $row;
        }
        echo empty($repos) ? "null" : json_encode($repos);
    }
    elseif ($key === 'commandspeak_history') {
        // Fetch ordered by timestamp descending
        $res = $db->query("SELECT * FROM history ORDER BY timestamp DESC");
        $hist = [];
        while ($row = $res->fetchArray(SQLITE3_ASSOC)) {
            $row['dryRun'] = (bool)$row['dryRun'];
            $hist[] = $row;
        }
        echo empty($hist) ? "null" : json_encode($hist);
    } else {
        echo "null";
    }
} 
elseif ($method === 'POST') {
    $input = json_decode(file_get_contents("php://input"), true);
    $key = $input['key'] ?? '';
    $value = $input['value'] ?? ''; // This comes in as a JSON string from JS
    $data = is_string($value) ? json_decode($value, true) : $value;
    
    if ($key === 'commandspeak_config') {
        $stmt = $db->prepare("REPLACE INTO config (key, value) VALUES ('user_config', :val)");
        $stmt->bindValue(':val', json_encode($data), SQLITE3_TEXT);
        $stmt->execute();
    }
    elseif ($key === 'commandspeak_repos') {
        $db->exec("DELETE FROM repositories"); // Clear old
        $stmt = $db->prepare("INSERT INTO repositories (id, name, path, active) VALUES (:id, :name, :path, :active)");
        foreach ($data as $repo) {
            $stmt->bindValue(':id', $repo['id'], SQLITE3_TEXT);
            $stmt->bindValue(':name', $repo['name'], SQLITE3_TEXT);
            $stmt->bindValue(':path', $repo['path'], SQLITE3_TEXT);
            $stmt->bindValue(':active', $repo['active'] ? 1 : 0, SQLITE3_INTEGER);
            $stmt->execute();
        }
    }
    elseif ($key === 'commandspeak_history') {
        $db->exec("DELETE FROM history"); // Clear old
        $stmt = $db->prepare("INSERT INTO history (id, repoId, query, intent, cmd, timestamp, dryRun) VALUES (:id, :repoId, :query, :intent, :cmd, :timestamp, :dryRun)");
        foreach ($data as $h) {
            $stmt->bindValue(':id', $h['id'], SQLITE3_TEXT);
            $stmt->bindValue(':repoId', $h['repoId'], SQLITE3_TEXT);
            $stmt->bindValue(':query', $h['query'], SQLITE3_TEXT);
            $stmt->bindValue(':intent', $h['intent'], SQLITE3_TEXT);
            $stmt->bindValue(':cmd', $h['cmd'], SQLITE3_TEXT);
            $stmt->bindValue(':timestamp', $h['timestamp'], SQLITE3_TEXT);
            $stmt->bindValue(':dryRun', $h['dryRun'] ? 1 : 0, SQLITE3_INTEGER);
            $stmt->execute();
        }
    }
    
    echo json_encode(["status" => "success"]);
}
?>
