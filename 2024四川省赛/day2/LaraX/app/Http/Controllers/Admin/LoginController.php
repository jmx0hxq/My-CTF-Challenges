<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use App\Models\Admin;
use Illuminate\Foundation\Auth\AuthenticatesUsers;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class LoginController extends Controller
{
    /*
    |--------------------------------------------------------------------------
    | Login Controller
    |--------------------------------------------------------------------------
    |
    | This controller handles authenticating users for the application and
    | redirecting them to your home screen. The controller uses a trait
    | to conveniently provide its functionality to your applications.
    |
    */

    use AuthenticatesUsers;

    public function __construct()
    {
        $this->middleware('guest')->except('logout');
    }

    public function showLoginForm()
    {
        return view('admin.login');
    }


    public function login(Request $request)
    {
        $this->validateLogin($request);
        var_dump(123);
        if ($this->hasTooManyLoginAttempts($request)) {
            $this->fireLockoutEvent($request);

            return $this->sendLockoutResponse($request);
        }

        if ($this->attemptLogin($request)) {
            return $this->sendLoginResponse($request);
        }

        $this->incrementLoginAttempts($request);

        return $this->sendFailedLoginResponse($request);
    }

    protected function attemptLogin(Request $request)
    {
        $credentials = [
            'email' => $request->input($this->username()),
            'password' => $request->input('password'),
        ];


        if ($status = $this->guard()->attempt($credentials)) {
            $user = $this->guard()->user();

            if($user->status_account != Admin::STATUS_ACCOUNT_VALID) {
                throw ValidationException::withMessages([
                    'account' => ['该帐户已被锁定，请联系客服。'],
                ]);
            }

            $user->ip_pre_login  = $user->ip_last_login;
            $user->ip_last_login = $request->getClientIp();

            $user->time_pre_login  = $user->time_last_login;
            $user->time_last_login = date('Y-m-d H:i:s');

            $user->count_login = $user->count_login + 1;

            $user->save();
        }

        return $status;
    }

    protected function redirectTo()
    {
        return '/admin/dashboard';
    }

    public function username()
    {
        return 'email';
    }

    protected function guard()
    {
        return Auth::guard('admin');
    }
}
