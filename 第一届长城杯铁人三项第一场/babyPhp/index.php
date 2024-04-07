<?php
//只保留了源码
function validate($text) {
    if (count(array_count_values(str_split($text))) > 7) {
        return false;
    }
    if (preg_match('/[a-z0-9]/i', $text)) {
        return false;
    }
    return true;
}

if (isset($_GET["exp"]) && !empty($_GET["exp"]) && is_string($_GET["exp"])) {
    echo $_GET["exp"];
    if(validate($_GET['exp']) === true) {
        eval($_GET['exp']);
    }
}else{
    highlight_file(__FILE__);
}
?>