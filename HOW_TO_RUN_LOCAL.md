# 🔒 Running CommandSpeak Locally with Private SQLite Database

Your data is **100% private** and stored on **your machine only**. It never goes to GitHub or any other server.

---

## How to Start (Windows)

### Step 1 — Make sure PHP is installed

Open PowerShell and check:
```powershell
php -v
```
If not installed, download from: https://windows.php.net/download/ (use the Thread Safe `.zip`, add to PATH).

### Step 2 — Start the local server

```powershell
cd C:\Users\devch\OneDrive\cyhi\commandspeak
php -S localhost:8080
```

You should see:
```
PHP Development Server (http://localhost:8080) started
```

### Step 3 — Open the app

Go to your browser and open:
```
http://localhost:8080/index.html
```

You'll see a green toast: **"🔒 Private SQLite DB active!"** — that means your data is being saved to `private_local_data.sqlite` on your local disk.

---

## What gets stored in the SQLite database

| Key | What it contains |
|-----|----------------|
| `commandspeak_config` | All your command templates, keywords, preferred deploy platform, package manager, theme |
| `commandspeak_repos` | List of repositories you've added |
| `commandspeak_history` | Every command you've translated/executed |

All of this is stored in `private_local_data.sqlite` in your project folder.

---

## Why is it private?

- `private_local_data.sqlite` is listed in `.gitignore` — it will **never** be pushed to GitHub
- The PHP server only listens on `localhost` — **no one else on the internet can access it**
- The `api.php` file itself is just a simple read/write handler with no authentication needed because only you can reach localhost

---

## Fallback (no PHP server)

If you open `index.html` directly as a file (`file://...`), the app automatically falls back to **browser localStorage**. Everything still works — just without the private SQLite backend.

---

## Stop the server

Just press `Ctrl+C` in the PowerShell window where PHP is running.
