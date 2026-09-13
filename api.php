<?php
header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');

// Determine home directory for Windows (USERPROFILE) or Linux/Mac (HOME)
$homeDir = getenv('USERPROFILE');
if (!$homeDir) {
    $homeDir = getenv('HOME');
}

$activityFile = $homeDir . DIRECTORY_SEPARATOR . '.commandspeak-activity.json';
$action = isset($_GET['action']) ? $_GET['action'] : 'activity';

if (!file_exists($activityFile)) {
    echo json_encode([]);
    exit;
}

$data = file_get_contents($activityFile);
$entries = json_decode($data, true);

if (!$entries) {
    echo json_encode([]);
    exit;
}

if ($action === 'repos') {
    $seen = [];
    $repos = [];
    foreach ($entries as $entry) {
        $name = isset($entry['repoName']) ? $entry['repoName'] : '';
        if ($name !== '' && !isset($seen[$name])) {
            $seen[$name] = true;
            $repos[] = [
                'name' => $name,
                'url' => isset($entry['repoUrl']) ? $entry['repoUrl'] : '',
                'path' => isset($entry['repoPath']) ? $entry['repoPath'] : ''
            ];
        }
    }
    echo json_encode($repos);
} else {
    // Default: return all activity
    echo json_encode($entries);
}
?>
