package com.oceancatch;

import android.app.Activity;
import android.os.Bundle;
import android.view.View;
import com.oceancatch.mobile.EbitenView;
import go.Seq;

public class MainActivity extends Activity {
    @Override public void onCreate(Bundle state) {
        super.onCreate(state);
        Seq.setContext(getApplicationContext());
        View view = new EbitenView(this);
        setContentView(view);
    }
}
