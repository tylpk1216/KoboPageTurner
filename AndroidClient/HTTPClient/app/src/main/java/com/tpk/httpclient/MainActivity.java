package com.tpk.httpclient;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.util.Log;
import android.view.KeyEvent;
import android.view.View;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.EditText;
import android.widget.MultiAutoCompleteTextView;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.StringRequest;
import com.android.volley.toolbox.Volley;

public class MainActivity extends AppCompatActivity {

    String mTag = "HTTPClient";

    RequestQueue mQueue;

    int mLeftCode = 21;
    int mRightCode = 22;

    MultiAutoCompleteTextView mNoteView;
    EditText mIPText;
    EditText mLeftText;
    EditText mRightText;

    Button mLeftPage;
    Button mRightPage;
    Button mExit;

    private void setUI() {
        String ip = PrefCtrl.getStringPref(this, "serverIP", "192.168.1.1");
        mIPText.setText(ip);

        String leftCode = PrefCtrl.getStringPref(this, "leftCode", String.valueOf(mLeftCode));
        mLeftText.setText(leftCode);
        mLeftCode = Integer.parseInt(leftCode);

        String rightCode = PrefCtrl.getStringPref(this, "rightCode", String.valueOf(mRightCode));
        mRightText.setText(rightCode);
        mRightCode = Integer.parseInt(rightCode);
    }

    private void recordUIPref() {
        String ip = mIPText.getText().toString().trim();
        PrefCtrl.recordStringPref(this, "serverIP", ip);

        String leftCode = mLeftText.getText().toString();
        PrefCtrl.recordStringPref(this, "leftCode", leftCode);

        String rightCode = mRightText.getText().toString();
        PrefCtrl.recordStringPref(this, "rightCode", rightCode);
    }

    @Override
    protected void onPause() {
        super.onPause();

        Log.d(mTag, "onPause");
        recordUIPref();
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();

        Log.d(mTag, "onDestroy");
        recordUIPref();
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);

        Log.d(mTag, "onCreate");

        mNoteView = (MultiAutoCompleteTextView) findViewById(R.id.NoteView);
        mIPText = (EditText) findViewById(R.id.edt_ip);
        mLeftText = (EditText) findViewById(R.id.edt_leftcode);
        mRightText = (EditText) findViewById(R.id.edt_rightcode);

        mLeftPage = (Button) findViewById(R.id.btn_left);
        mRightPage = (Button) findViewById(R.id.bth_right);
        mExit = (Button) findViewById(R.id.btn_exit);

        mQueue = Volley.newRequestQueue(this.getApplicationContext());

        setUI();
    }

    public void disableButton() {
        mLeftPage.setClickable(false);
        mRightPage.setClickable(false);
        mExit.setClickable(false);
    }

    public void enableButton() {
        mLeftPage.setClickable(true);
        mRightPage.setClickable(true);
        mExit.setClickable(true);
    }

    public void clearLog() {
        mIPText.clearFocus();
        mLeftText.clearFocus();
        mRightText.clearFocus();
        mNoteView.setText("");
    }

    public void sendHTTPRequest(String url) {
        Log.d(mTag, url);

        clearLog();
        disableButton();

        StringRequest stringRequest = new StringRequest(Request.Method.GET, url,
                new Response.Listener<String>() {
                    @Override
                    public void onResponse(String response) {
                        String t = response.substring(0, 19);
                        int i = response.indexOf(") (");
                        String res = "";
                        if (i != -1) {
                            res = response.substring(i + 3).replace(")", "");
                        }

                        res = t + " " +  res;
                        mNoteView.setText(res);
                        enableButton();
                    }
                }, new Response.ErrorListener() {
                @Override
                public void onErrorResponse(VolleyError error) {
                    mNoteView.setText(error.toString());
                    enableButton();
                }
            }
        );

        // Add the request to the RequestQueue.
        mQueue.add(stringRequest);
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event) {
        String sCode = String.valueOf(keyCode);

        Log.d(mTag, sCode);
        mNoteView.setText("Your key code is " + sCode);

        if (mLeftText.isFocused()) {
            mLeftCode = keyCode;
            mLeftText.setText(sCode);
            recordUIPref();
        } else if (mRightText.isFocused()) {
            mRightCode = keyCode;
            mRightText.setText(sCode);
            recordUIPref();
        } else {
            if (keyCode == mLeftCode) {
                sendAPI("/left");
            } else if (keyCode == mRightCode) {
                sendAPI("/right");
            }
        }

        //return super.onKeyDown(keyCode, event);
        return true;
    }

    public void sendAPI(String endpoint) {
        String url = "http://" + mIPText.getText().toString() + endpoint;
        sendHTTPRequest(url);
    }

    public void leftPage(View view) {
        Log.d(mTag, "leftPage");

        sendAPI("/left");
    }

    public void rightPage(View view) {
        Log.d(mTag, "rightPage");

        sendAPI("/right");
    }

    public void exitServer(View view) {
        Log.d(mTag, "exitServer");

        sendAPI("/exit");
        recordUIPref();
    }
}

