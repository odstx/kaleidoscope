-addkeepattributes
-keepattributes Signature
-keepattributes *Annotation*
-keep class com.example.app.data.model.** { *; }
-keep class retrofit2.** { *; }
-keepclasseswithmembers class * {
    @retrofit2.http.* <methods>;
}
