<?php


function _loader($class)
{
    $class_file = './class/' . (str_replace('\\', '/', $class));
    if(file_exists($class_file)){
        include './class/' . (str_replace('\\', '/', $class));
    }else {
        include './class/' . (str_replace('\\', '/', $class) . '.php');
    }
}

spl_autoload_register("_loader");




