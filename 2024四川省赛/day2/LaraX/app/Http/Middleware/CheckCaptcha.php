<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Support\Facades\Session;
use Illuminate\Validation\ValidationException;

class CheckCaptcha
{
    public function handle($request, Closure $next)
    {
        
        Session::remove('captcha');

        return $next($request);
    }
}
