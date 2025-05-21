const fs = require('fs');
const path = require('path');

// Helper function to generate UUID
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// Parse SRT time to microseconds
function parseSRTTime(timeStr) {
    const [timePart, millisPart] = timeStr.split(',');
    const [hours, minutes, seconds] = timePart.split(':').map(Number);
    const millis = parseInt(millisPart, 10);
    return hours * 3600000000 + minutes * 60000000 + seconds * 1000000 + millis * 1000;
}

// Parse SRT file
function parseSRT(filename) {
    const content = fs.readFileSync(filename, 'utf-8');
    const lines = content.split('\n');
    const subtitles = [];
    let currentSubtitle = null;

    for (const line of lines) {
        const trimmed = line.trim();
        if (!trimmed) continue;

        // Check if line is a number (subtitle index)
        if (/^\d+$/.test(trimmed)) {
            if (currentSubtitle) {
                subtitles.push(currentSubtitle);
            }
            currentSubtitle = { Index: subtitles.length + 1 };
            continue;
        }

        // Check if line is a time range
        const timeMatch = trimmed.match(/(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})/);
        if (timeMatch && currentSubtitle) {
            currentSubtitle.Start = timeMatch[1];
            currentSubtitle.End = timeMatch[2];
            const startTime = parseSRTTime(timeMatch[1]);
            const endTime = parseSRTTime(timeMatch[2]);
            currentSubtitle.Duration = endTime - startTime;
            continue;
        }

        // Otherwise, it's subtitle text
        if (currentSubtitle) {
            currentSubtitle.Text = currentSubtitle.Text 
                ? currentSubtitle.Text + '\n' + trimmed 
                : trimmed;
        }
    }

    if (currentSubtitle) {
        subtitles.push(currentSubtitle);
    }

    return subtitles;
}

// Convert subtitles to JSON format
function convertToJSON(subtitles) {
    const output = {
        canvas_config: {
            height: 1080,
            ratio: "original",
            width: 1920
        },
        color_space: 0,
        config: {
            adjust_max_index: 1,
            attachment_info: [],
            combination_max_index: 1,
            export_range: null,
            extract_audio_last_index: 1,
            lyrics_recognition_id: "",
            lyrics_sync: true,
            lyrics_taskinfo: [],
            maintrack_adsorb: true,
            material_save_mode: 0,
            original_sound_last_index: 1,
            record_audio_last_index: 1,
            sticker_max_index: 1,
            subtitle_recognition_id: "",
            subtitle_sync: true,
            subtitle_taskinfo: [{
                id: generateUUID(),
                language: "",
                remove_invalid_task_id: "",
                type: 10
            }],
            video_mute: false,
            zoom_info_params: null
        },
        cover: null,
        create_time: 0,
        duration: 32600000,
        extra_info: null,
        fps: 30.0,
        free_render_index_mode_on: false,
        group_container: null,
        id: generateUUID(),
        keyframes: {
            adjusts: [],
            audios: [],
            filters: [],
            handwrites: [],
            stickers: [],
            texts: [],
            videos: []
        },
        last_modified_platform: {
            app_id: 359289,
            app_source: "cc",
            app_version: "1.5.0",
            device_id: "839a3a0281bf298bb7a04ef106f6f838",
            hard_disk_id: "2042ebf3be3c78787b62a7cf8ea27d5d",
            mac_address: "72859f122da5bc1c1d727bfd8490ee4f",
            os: "windows",
            os_version: "10.0.19044"
        },
        materials: {
            audio_balances: [],
            audio_effects: [],
            audio_fades: [],
            audios: [],
            beats: [],
            canvases: [],
            chromas: [],
            color_curves: [],
            drafts: [],
            effects: [],
            handwrites: [],
            hsl: [],
            images: [],
            log_color_wheels: [],
            manual_deformations: [],
            masks: [],
            material_animations: [
                {
                    animations: [],
                    id: generateUUID(),
                    type: "sticker_animation"
                },
                {
                    animations: [],
                    id: generateUUID(),
                    type: "sticker_animation"
                }
            ],
            placeholders: [],
            plugin_effects: [],
            primary_color_wheels: [],
            realtime_denoises: [],
            speeds: [],
            stickers: [],
            tail_leaders: [],
            text_templates: [],
            texts: [],
            transitions: [],
            video_effects: [],
            video_trackings: [],
            videos: []
        },
        mutable_config: null,
        name: "",
        new_version: "68.0.1",
        platform: {
            app_id: 359289,
            app_source: "cc",
            app_version: "1.5.0",
            device_id: "839a3a0281bf298bb7a04ef106f6f838",
            hard_disk_id: "2042ebf3be3c78787b62a7cf8ea27d5d",
            mac_address: "72859f122da5bc1c1d727bfd8490ee4f",
            os: "windows",
            os_version: "10.0.19044"
        },
        relationships: [],
        render_index_track_mode_on: false,
        retouch_cover: null,
        source: "default",
        static_cover_image_path: "",
        tracks: [
            {
                attribute: 0,
                flag: 0,
                id: generateUUID(),
                segments: [],
                type: "video"
            },
            {
                attribute: 0,
                flag: 1,
                id: generateUUID(),
                segments: [],
                type: "text"
            }
        ],
        update_time: 0,
        version: 360000
    };

    // Add subtitles to the materials.texts and tracks
    subtitles.forEach((sub, i) => {
        const startTime = parseSRTTime(sub.Start);
        const textId = generateUUID();
        const animationId = output.materials.material_animations[i % 2].id;

        // Add text entry
        output.materials.texts.push({
            add_type: 1,
            alignment: 1,
            background_alpha: 1.0,
            background_color: "",
            background_height: 1.0,
            background_horizontal_offset: 0.0,
            background_round_radius: 0.0,
            background_style: 0,
            background_vertical_offset: 0.0,
            background_width: 1.0,
            bold_width: 0.0,
            border_color: "",
            border_width: 0.08,
            check_flag: 7,
            content: `<font id="" path="C:/Users/os/AppData/Local/CapCut/Apps/1.5.0.230/Resources/Font/SystemFont/en.ttf"><color=(1.000000, 1.000000, 1.000000, 1.000000)><size=5.000000>[${sub.Text}]</size></color></font>`,
            font_category_id: "",
            font_category_name: "",
            font_id: "",
            font_name: "",
            font_path: "C:/Users/os/AppData/Local/CapCut/Apps/1.5.0.230/Resources/Font/SystemFont/en.ttf",
            font_resource_id: "",
            font_size: 5.0,
            font_source_platform: 0,
            font_team_id: "",
            font_title: "none",
            font_url: "",
            fonts: [],
            global_alpha: 1.0,
            group_id: "",
            has_shadow: false,
            id: textId,
            initial_scale: 1.0,
            is_rich_text: false,
            italic_degree: 0,
            ktv_color: "",
            layer_weight: 1,
            letter_spacing: 0.0,
            line_spacing: 0.02,
            name: "",
            recognize_type: 0,
            shadow_alpha: 0.8,
            shadow_angle: -45.0,
            shadow_color: "",
            shadow_distance: 8.0,
            shadow_point: {
                x: 1.0182337649086284,
                y: -1.0182337649086284
            },
            shadow_smoothing: 1.0,
            shape_clip_x: false,
            shape_clip_y: false,
            style_name: "",
            sub_type: 0,
            text_alpha: 1.0,
            text_color: "#FFFFFF",
            text_preset_resource_id: "",
            text_size: 30,
            text_to_audio_ids: [],
            tts_auto_update: false,
            type: "subtitle",
            typesetting: 0,
            underline: false,
            underline_offset: 0.22,
            underline_width: 0.05,
            use_effect_default_color: true,
            words: []
        });

        // Add segment to text track
        const segment = {
            cartoon: false,
            clip: {
                alpha: 1.0,
                flip: {
                    horizontal: false,
                    vertical: false
                },
                rotation: 0.0,
                scale: {
                    x: 1.0,
                    y: 1.0
                },
                transform: {
                    x: 0.0,
                    y: -0.73
                }
            },
            enable_adjust: false,
            enable_color_curves: true,
            enable_color_wheels: true,
            enable_lut: false,
            extra_material_refs: [animationId],
            group_id: "",
            hdr_settings: null,
            id: generateUUID(),
            intensifies_audio: false,
            is_placeholder: false,
            is_tone_modify: false,
            keyframe_refs: [],
            last_nonzero_volume: 1.0,
            material_id: textId,
            render_index: 14000 - i,
            reverse: false,
            source_timerange: null,
            speed: 1.0,
            target_timerange: {
                duration: sub.Duration,
                start: startTime
            },
            template_id: "",
            track_attribute: 0,
            track_render_index: 0,
            visible: true,
            volume: 1.0
        };

        // Find the text track and add the segment
        const textTrack = output.tracks.find(t => t.type === "text");
        if (textTrack) {
            textTrack.segments.push(segment);
        }
    });

    return output;
}

// Main function
function main() {
    const inputFile = "subtitles.srt";
    const outputFile = "C:/Users/os/AppData/Local/CapCut/User Data/Projects/com.lveditor.draft/abc/draft_content.json";

    try {
        // Check if input file exists
        if (!fs.existsSync(inputFile)) {
            console.error(`Error: ${inputFile} not found in current directory`);
            process.exit(1);
        }

        // Parse and convert
        const subtitles = parseSRT(inputFile);
        const jsonOutput = convertToJSON(subtitles);
        
        // Write output
        fs.writeFileSync(outputFile, JSON.stringify(jsonOutput));
        console.log(`Successfully converted ${inputFile} to ${outputFile}`);
    } catch (err) {
        console.error(`Error: ${err.message}`);
        process.exit(1);
    }
}

// Run the program
main();