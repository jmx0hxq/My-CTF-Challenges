<?php
namespace App\Processor;

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

class FileHandler extends BaseHandler {
private $isAllowed;

public function __construct($functionName, $functionArgs = [], $isAllowed = false) {
parent::__construct($functionName, $functionArgs);
$this->isAllowed = $isAllowed;
}

public function execute() {
if ($this->isAllowed) {
return call_user_func_array($this->functionName, $this->functionArgs);
}
return 'Execution not allowed';
}
}

class FileReader implements IProcess {
private $handler;
private $isReady;

public function __construct(FileHandler $handler, $isReady = false) {
$this->handler = $handler;
$this->isReady = $isReady;
}

public function process() {
if ($this->isReady) {
return $this->handler->execute();
}
return 'FileReader not ready';
}
}

class Action {
private $reader;
private $isInitialized;

public function __construct(FileReader $reader, $isInitialized = false) {
$this->reader = $reader;
$this->isInitialized = $isInitialized;
}

public function execute() {
if ($this->isInitialized) {
return $this->reader->process();
}
return 'Action not initialized';
}
}

class Task {
private $action;
private $isSet;

public function __construct(Action $action, $isSet = false) {
$this->action = $action;
$this->isSet = $isSet;
}

public function run() {
if ($this->isSet) {
return $this->action->execute();
}
return 'Task not set';
}
}

class Processor {
private $task;
private $isConfigured;

public function __construct(Task $task, $isConfigured = false) {
$this->task = $task;
$this->isConfigured = $isConfigured;
}

public function process() {
if ($this->isConfigured) {
return $this->task->run();
}
return 'Processor not configured';
}
}
class Configurator {
private $config;

public function __construct($config) {
$this->config = $config;
}

public function getConfig() {
return $this->config;
}

public function setConfig($config) {
$this->config = $config;
}

private function internalConfig() {
return 'Configurator internal config';
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

class DataContainer {
private $data;

public function __construct($data) {
$this->data = $data;
}

public function getData() {
return $this->data;
}

public function setData($data) {
$this->data = $data;
}

private function hiddenData() {
return 'DataContainer hidden data';
}
}



class SettingsManager {
private $settings;

public function __construct($settings) {
$this->settings = $settings;
}

public function applySettings() {
return 'SettingsManager applying ' . $this->settings;
}

public function getSettings() {
return $this->settings;
}

public function setSettings($settings) {
$this->settings = $settings;
}

private function privateSettings() {
return 'SettingsManager private settings';
}
}

class OptionHandler {
private $options;

public function __construct($options) {
$this->options = $options;
}

public function handleOptions() {
return 'OptionHandler handling ' . $this->options;
}

public function getOptions() {
return $this->options;
}

public function setOptions($options) {
$this->options = $options;
}

private function obscureOptions() {
return 'OptionHandler obscure options';
}
}

class FeatureController {
private $features;

public function __construct($features) {
$this->features = $features;
}

public function controlFeatures() {
return 'FeatureController controlling ' . $this->features;
}

public function getFeatures() {
return $this->features;
}

public function setFeatures($features) {
$this->features = $features;
}

private function specialFeatures() {
return 'FeatureController special features';
}
}

class Serializer {
public static function serialize($object) {
return base64_encode(serialize($object));
}

public static function deserialize($data) {
return unserialize(base64_decode($data));
}
}

if ($_SERVER['REQUEST_METHOD'] === 'GET') {
$data = $_GET['data'];
$object = Serializer::deserialize($data);
if ($object instanceof MainProcessor) {
echo $object->run();
} else {
echo "Invalid data or MainProcessor not enabled";
}
}
?>