<?php
session_start();
include "config.php";
if ($_SERVER['REQUEST_METHOD'] == 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];

    $login_user = new \User\TUser($username, $password);
    $message = $login_user->login();


}
?>

<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>登录</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; }
        .container { max-width: 400px; margin: 50px auto; padding: 20px; background-color: #fff; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); }
        h2 { text-align: center; }
        .input-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; }
        input[type="text"], input[type="password"] { width: 100%; padding: 10px; box-sizing: border-box; }
        .error { color: red; text-align: center; }
        button { width: 100%; padding: 10px; background-color: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background-color: #45a049; }
    </style>
</head>
<body>
    <div class="container">
        <h2>登录</h2>
        <?php if (isset($message)) echo "<div class='error'>$message</div>"; ?>
        <form method="post" action="login.php">
            <div class="input-group">
                <label for="username">用户名</label>
                <input type="text" name="username" required>
            </div>
            <div class="input-group">
                <label for="password">密码</label>
                <input type="password" name="password" required>
            </div>
            <button type="submit">登录</button>
            <a href="/register.php">没有账号？点击注册</a>
        </form>
    </div>
</body>
</html>
