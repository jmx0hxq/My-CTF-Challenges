<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>
            SuperAdmin
        </title>
        <meta name="renderer" content="webkit">
        <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
        <meta name="apple-mobile-web-app-status-bar-style" content="black">
        <meta name="apple-mobile-web-app-capable" content="yes">
        <meta name="format-detection" content="telephone=no">
        <link rel="stylesheet" href="./css/x-admin.css" media="all">
    </head>
    <body>
        <div class="x-nav">
            <span class="layui-breadcrumb">
              <a><cite>首页</cite></a>
              <a><cite>会员管理</cite></a>
              <a><cite>系统日志</cite></a>
            </span>
            <a class="layui-btn layui-btn-small" style="line-height:1.6em;margin-top:3px;float:right"  href="javascript:location.replace(location.href);" title="刷新"><i class="layui-icon" style="line-height:30px">ဂ</i></a>
        </div>
        <div class="x-body">
            <table class="layui-table">
                <tbody>
                <?php
                include "api/conn.php";
                foreach ($pdo->query('SELECT * from logs') as $row) {
                    $str = <<<EOF
                    <tr>
                        <td>
                            CONTENT
                        </td>
                    </tr>
EOF;
                    $str = str_replace("CONTENT",htmlspecialchars($row[1]),$str);
                    echo $str;
                    if(!preg_match("/php/i", $str)) {
                        file_put_contents("temp/.temp", htmlspecialchars_decode($str));
                        include "temp/.temp";
                        unlink("temp/.temp");
                    }
                }
                ?>
                </tbody>
            </table>
        </div>
        <script src="./lib/layui/layui.js" charset="utf-8"></script>
        <script src="./js/x-layui.js" charset="utf-8"></script>
    </body>
</html>
