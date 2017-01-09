# -*- coding: utf-8 -*-
# Generated by Django 1.10.4 on 2016-12-20 12:02
from __future__ import unicode_literals

import django.contrib.gis.db.models.fields
from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Scan',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('scan_id', models.IntegerField()),
                ('scan_moment', models.DateField()),
                ('scan_source', models.CharField(max_length=15)),
                ('afstand', models.CharField(max_length=25, null=True)),
                ('latitude', models.DecimalField(decimal_places=8, max_digits=13)),
                ('longitude', models.DecimalField(decimal_places=8, max_digits=13)),
                ('buurtcode', models.CharField(max_length=4, null=True)),
                ('sperscode', models.CharField(max_length=15)),
                ('qualcode', models.CharField(max_length=35, null=True)),
                ('ff_df', models.CharField(max_length=15, null=True)),
                ('nha_nr', models.IntegerField(null=True)),
                ('nha_hoogte', models.DecimalField(decimal_places=3, max_digits=6, null=True)),
                ('uitval_nachtrun', models.CharField(max_length=8, null=True)),
                ('geometrie', django.contrib.gis.db.models.fields.PointField(null=True, srid=4326)),
                ('geometrie_rd', django.contrib.gis.db.models.fields.PointField(null=True, srid=28992)),
            ],
        ),
    ]
