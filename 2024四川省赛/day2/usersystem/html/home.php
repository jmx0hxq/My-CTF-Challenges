<?php
session_start();

if (!isset($_SESSION['username'])) {
    header("Location: login.php");
    exit;
}

$username = $_SESSION['username'];

?>

<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>主页</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; }
        .container { max-width: 400px; margin: 50px auto; padding: 20px; background-color: #fff; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: center; }
        h2 { margin-bottom: 20px; }
        a { display: inline-block; margin: 10px 0; color: #4CAF50; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="container">
        <h2>欢迎, <?php echo htmlspecialchars($username);?></h2>
        <p><a href="change_password.php">修改密码</a></p>
        <p><a href="logout.php">登出</a></p>
    </div>
</body>
</html>
