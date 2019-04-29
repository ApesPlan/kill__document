<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateActivityAlertStyle extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('activity_alert_style', function (Blueprint $table) {
            $table->increments('id', 11)->comment('弹窗式活动id');
            $table->string('activity_name', 100)->default('')->comment('活动名称');
            $table->timestamp('begin_time')->comment('开始时间');
            $table->timestamp('end_time')->comment('结束时间');
            $table->string('img', 255)->default('')->comment('海报图片');
            $table->unsignedSmallInteger('width', 5)->default(0)->comment('弹窗宽度');
            $table->unsignedSmallInteger('hight', 5)->default(0)->comment('弹窗高度');
            $table->unsignedSmallInteger('top', 5)->default(0)->comment('弹窗与顶部的距离');
            $table->unsignedSmallInteger('bottom', 5)->default(0)->comment('弹窗与底部侧的距离');
            $table->unsignedSmallInteger('left', 5)->default(0)->comment('弹窗与左侧的距离');
            $table->unsignedSmallInteger('right', 5)->default(0)->comment('弹窗与右侧的距离');
            $table->unsignedTinyInteger('position_type', 1)->default(0)->comment('弹窗位置类型：左上，左下，右下，右上，居中');
            $table->unsignedTinyInteger('status', 1)->default(0)->comment('弹窗状态');
            $table->string('target_url', 255)->default('')->comment('弹窗跳转地址');
            $table->text('display_where')->default('')->comment('图片展示页面，json字符串');
            $table->timestamps();
            $table->engine = 'InnoDB';
            $table->charset = 'utf8';
            $table->collation = 'utf8_unicode_ci';
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('activity_alert_style');
    }
}
