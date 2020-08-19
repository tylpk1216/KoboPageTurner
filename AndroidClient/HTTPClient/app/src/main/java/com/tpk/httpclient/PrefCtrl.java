package com.tpk.httpclient;

import android.content.Context;
import android.content.SharedPreferences;
import android.content.SharedPreferences.Editor;

public class PrefCtrl {

    private static String mKeyFile = "HTTPClient";

    public PrefCtrl() {
        // TODO Auto-generated constructor stub
    }

    public static void recordStringPref(Context context, String prefKey, String prefValue) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        Editor editor = settings.edit();
        editor.putString(prefKey, prefValue);
        editor.commit();
    }

    public static String getStringPref(Context context, String prefKey, String defaultVal) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        if(settings != null) {
            return settings.getString(prefKey, defaultVal);
        }

        return defaultVal;
    }

    public static void recordIntPref(Context context, String prefKey, int prefValue) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        Editor editor = settings.edit();
        editor.putInt(prefKey, prefValue);
        editor.commit();
    }

    public static int getIntPref(Context context, String prefKey, int defaultVal) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        if(settings != null) {
            return settings.getInt(prefKey, defaultVal);
        }

        return defaultVal;
    }

    public static void recordBoolPref(Context context,String prefKey, Boolean prefValue) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        Editor editor = settings.edit();
        editor.putBoolean(prefKey, prefValue);
        editor.commit();
    }

    public static Boolean getBoolPref(Context context,String prefKey, Boolean defaultVal) {
        SharedPreferences settings = context.getSharedPreferences(mKeyFile, Context.MODE_PRIVATE);
        if(settings != null) {
            return settings.getBoolean(prefKey, defaultVal);
        }

        return defaultVal;
    }

}
