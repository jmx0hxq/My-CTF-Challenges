<?php
session_start();
include "config.php";
if (!isset($_SESSION['username'])) {
    header("Location: login.php");
    exit;
}



$username = $_SESSION['username'];


if ($_SERVER['REQUEST_METHOD'] == 'POST') {
    $username = $_POST['username'];
    $new_password = $_POST['new_password'];
    $confirm_password = $_POST['confirm_password'];

    if ($new_password !== $confirm_password) {
        $message = "两次输入的密码不一致";
    } else {
        $change_user = new \User\TUser($username, $new_password);
        $message = $change_user->change_password();

    }
}
?>

<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>修改密码</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; }
        .container { max-width: 400px; margin: 50px auto; padding: 20px; background-color: #fff; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); }
        h2 { text-align: center; }
        .input-group { margin: 15px 0; }
        label { display: block; margin-bottom: 5px; }
        input[type="password"] { width: 100%; padding: 10px; box-sizing: border-box; }
        .error { color: red; text-align: center; }
        .success { color: green; text-align: center; }
        button { width: 100%; padding: 10px; background-color: #4CAF50; color: white; border: none; cursor: pointer; }
        button:hover { background-color: #45a049; }
    </style>
</head>
<body>
    <div class="container">
        <h2>修改密码</h2>
        <?php if (isset($message)) echo "<div class='success'>$message</div>"; ?>
        <form method="post" action="change_password.php">
        <input type="hidden" name="username" value="<?php echo htmlspecialchars($username); ?>">
            <div class="input-group">
                <label for="new_password">新密码</label>
                <input type="password" name="new_password" required>
            </div>
            <div class="input-group">
                <label for="confirm_password">确认新密码</label>
                <input type="password" name="confirm_password" required>
            </div>
            <button type="submit">修改密码</button>
        </form>
    </div>
</body>
</html>
