<?php

namespace Safe\Rule{

    class Check{
        public $callback;

        public $username;

        public function __construct($username){

            $this->username = $username;


        }

        function isUsernameValid($username): bool
        {
            if (strpos($username, '.ph') !== false) {
                return false;
            }
            return true;
        }

        public function __invoke(): bool
        {

            return $this->isUsernameValid($this->username);

        }


    }
}
