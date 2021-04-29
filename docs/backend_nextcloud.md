# Using Nextcloud as a Growlapse Backend

Nextcloud can be used as a backend for Growlapse, but for safe access, some steps are required. Let's go through them one by one! Replace `$NX_URL` with your Nextcloud URL in the following, i.e. with `https://asdf.your-storageshare.de`.

## Set up Nextcloud for the Agent

1. Login as the admin user
2. Visit `$NX_URL/settings/users`
3. Click on "New user", use `growlapse-agent` as the username and set a password (I'll be using `agentpass` in the following), and create it

## Connect the Agent to Nextcloud

1. Follow the installation and usage instructions at https://github.com/pojntfx/growlapse#installation
2. Use the following parameters when starting `growlapse-agent`: `-webdavURL $NX_URL/remote.php/dav/files/growlapse-agent/ -webdavUsername growlapse-agent -webdavPassword agentpass`

## Set up Nextcloud for the Frontend

1. Install [WebAppPassword](https://apps.nextcloud.com/apps/webapppassword) on your Nextcloud
2. Visit `$NX_URL/settings/admin/webapppassword` and add `https://pojntfx.github.io` to "Allowed origins"
3. Install [Guests](https://apps.nextcloud.com/apps/guests) on your Nextcloud
4. Visit `$NX_URL/settings/admin/guests` and add `webapppassword` to "app whitelist"
5. Logout
6. Login as `growlapse-agent`
7. Share the folder you've chosen as the `-webdavPrefix` (`/Growlapse` by default, if it doesn't exist, you can create it now) with an email address you have access to (this will be the guest user's email, I'll use `growlapse-frontend@example.com` here), then click on **Invite guest**
8. Accept the invitation sent to the mail address and set a password (I'll be using `frontendpass` in the following; this user is read-only so you can safely share it to those who should be able to access your Growlapse data)

## Connect the Frontend to Nextcloud

1. Visit https://pojntfx.github.io/growlapse/
2. Set `$NX_URL/remote.php/dav/files/growlapse-frontend@example.com/` as the API URL
3. Set `growlapse-frontend@example.com` as the username
4. Set `frontendpass` as the password
5. Click on "Login"

That's it! You've now configured Growlapse to safely use Nextcloud as a backend. To share your Growlapse data, simply share the URL, which is pre-authenticated and should look something like this: https://pojntfx.github.io/growlapse/?webDAVPassword=agentpass&webDAVURL=https%3A%2F%2Fexample.your-storageshare.de%2Fremote.php%2Fdav%2Ffiles%2Fgrowlapse-frontend%40example.com%2F&webDAVUsername=growlapse-frontend%40example.com - you can also bookmark it to have quick access to it.
