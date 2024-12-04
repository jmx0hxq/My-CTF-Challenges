<?php
session_start();
include "config.php";
if ($_SERVER['REQUEST_METHOD'] == 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $confirm_password = $_POST['confirm_password'];
    if ($password !== $confirm_password) {
        $message = "两次输入的密码不一致";

    }else{
        $register_user = new \User\TUser(basename($username), $password);
        $message = $register_user->register();
    }

}
?>

<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>注册</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; }
        .container { max-width: 400px; margin: 50px auto; padding: 20px; background-color: #fff; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); }
        h2 { text-align: center; }
        .input-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; }
        input[type="text"], input[type="password"] { width: 100%; padding: 10px; box-sizing: border-box; }
        .error { color: red; text-align: center; }
        .success { color: green; text-align: center; }
        button { width: 100%; padding: 10px; background-color: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background-color: #45a049; }
    </style>
</head>
<body>
    <div class="container">
        <h2>注册</h2>
        <?php if (isset($message)) echo "<div>$message</div>"; ?>

        <form method="post" action="register.php">
            <div class="input-group">
                <label for="username">用户名</label>
                <input type="text" name="username" required>
            </div>
            <div class="input-group">
                <label for="password">密码</label>
                <input type="password" name="password" required>
            </div>
            <div class="input-group">
                <label for="confirm_password">确认密码</label>
                <input type="password" name="confirm_password" required>
            </div>
            <a href="/login.php">已有账号？点击登陆</a>
            <button type="submit">注册</button>
        </form>
    </div>
</body>
</html>
