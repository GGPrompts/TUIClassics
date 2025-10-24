# Supabase Setup Guide

Quick reference for setting up the Supabase backend for TUI Arcade.

---

## Step 1: Create Supabase Account (2 minutes)

1. Go to https://supabase.com
2. Click "Start your project"
3. Sign up with GitHub (recommended) or email
4. Verify email if needed

---

## Step 2: Create New Project (3 minutes)

1. Click "New Project"
2. Fill in:
   - **Name**: `tuiarcade`
   - **Database Password**: (generate strong password - save it!)
   - **Region**: Choose closest to you
   - **Pricing Plan**: Free

3. Click "Create new project"
4. Wait 2-3 minutes for provisioning

---

## Step 3: Get API Keys (1 minute)

1. In Supabase dashboard, click **Settings** (gear icon in sidebar)
2. Click **API** in settings menu
3. Copy these two values:

```bash
Project URL:
https://xxxxxxxxxxxxx.supabase.co

anon public key:
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

4. Save them in `~/projects/TUIClassics/.env`:

```bash
# Create .env file
cat > ~/projects/TUIClassics/.env << 'EOF'
SUPABASE_URL=https://YOUR_PROJECT_ID.supabase.co
SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
EOF
```

**Important**: Add `.env` to `.gitignore`:
```bash
echo ".env" >> ~/projects/TUIClassics/.gitignore
```

---

## Step 4: Run Database Schema (5 minutes)

1. In Supabase dashboard, click **SQL Editor** in sidebar
2. Click **New Query**
3. Copy the entire SQL schema from `MULTIPLAYER.md` Section 1.2
4. Paste into the SQL editor
5. Click **Run** (or press Ctrl+Enter)

You should see: ✅ Success. No rows returned

---

## Step 5: Verify Tables Created

1. Click **Table Editor** in sidebar
2. You should see these tables:
   - `players`
   - `scores`
   - `games`
   - `game_participants`
   - `chat_messages`
   - `achievements`

3. Click on `players` table - it should be empty (for now)

---

## Step 6: Test Insert (Optional)

To verify everything works, try inserting a test player:

1. In SQL Editor, run:

```sql
insert into players (username) values ('TestPlayer');
select * from players;
```

2. You should see the test player with a UUID

3. Delete test data:

```sql
delete from players where username = 'TestPlayer';
```

---

## Step 7: Ready to Code!

Your Supabase backend is ready. Return to your terminal and tell Claude:

```
"Create the internal/supabase package with client initialization"
```

---

## Troubleshooting

### "relation does not exist" error
- You forgot to run the schema SQL
- Go back to Step 4

### "unauthorized" when querying
- Check your anon key is correct in `.env`
- Verify RLS policies are set up (they're in the schema)

### Can't connect from Go
- Verify `SUPABASE_URL` has no trailing slash
- Check you installed: `go get github.com/supabase-community/supabase-go`

### Tables not showing in Table Editor
- Refresh the page
- Check SQL ran successfully (no red errors)

---

## Quick Reference

**Dashboard URL**: https://supabase.com/dashboard/project/YOUR_PROJECT_ID

**Direct Links:**
- Table Editor: `/editor`
- SQL Editor: `/sql`
- API Docs: `/api`
- Settings: `/settings/api`

---

**Setup complete! Time to build. 🚀**
