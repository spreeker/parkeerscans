# -*- coding: utf-8 -*-
# Generated by Django 1.10.4 on 2016-12-20 12:23
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('scans', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='scan',
            name='scan_moment',
            field=models.DateTimeField(),
        ),
    ]
